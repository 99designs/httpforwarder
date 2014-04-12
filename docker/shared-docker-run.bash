# Wrapper for `docker run` to handle volume mounting and container linking.
# Run with --help for usage, or read source.
#
# Designed to be sourced from a per-project docker-run.sh which has set the
# following variables:

# CONTAINER_NAME: The unique name of the container.
# DOCKER_ENVS: Env options, e.g. "--env DB_HOST=localhost --env DB_USER=root"
# DOCKER_IMAGE: Docker image to run.
# DOCKER_LINKS: Link options, e.g. "--link mysql-5.6:mysql"
# DOCKER_OPTIONAL_LINKS: Optional links, e.g. "redis mysql-5.6:mysql"
# DOCKER_PORTS: Port options, e.g. "--publish 49001:80"
# DOCKER_VOLUMES: Volume options, e.g. "--volume /projects/commerce:/commerce"
# MY_INIT: Non-empty to run command via my_init

##
# Functions

usage() {
  cat <<END
Usage: $0 [-i|-r] [command]
  -i: run with --interactive --tty --rm, default command to bash.
  -r: run with --rm (requires a command; incompatible with detached mode)
  If no command nor -i, runs --detach with default Dockerfile command.

END
}
if [ "$1" = "-h" -o "$1" = "--help" ]; then usage; exit 0; fi

function is_container_running() {
  docker inspect --format='{{.State.Running}}' $1 2>&1 | grep -q true
}

##
# Argument parsing.

# DOCKER_CMD is an array, not a string.
# This is to preserve quoted grouping of space-separated args.
# e.g. ./docker-run.sh bash -c 'date "+%Y %M %d"'
#      $1: "bash", $2: "-c", $3: "date \"+%Y %M %d\""
DOCKER_CMD=()

while getopts "ir" opt; do
  case "$opt" in
    "i")
      DOCKER_OPTS+=" --interactive --tty --rm"
      DOCKER_CMD=("bash")
      ;;
    "r")
      DOCKER_OPTS+=" --rm"
      ;;
  esac
done
shift $(( $OPTIND - 1))

# $@ (if present) takes precedence over an inferred DOCKER_CMD
if [ $# -gt 0 ]; then
  DOCKER_CMD=("$@")
fi

# If no command, run detached, named, with ports exposed.
if [ -z "$DOCKER_CMD" ]; then
  DOCKER_OPTS="$DOCKER_OPTS --detach --name $CONTAINER_NAME $DOCKER_PORTS"
  delete_old_container="delete_old_container"
fi

if [ -n "$MY_INIT" -a -n "$DOCKER_CMD" ]; then
  DOCKER_CMD_PREFIX="/sbin/my_init --quiet --skip-runit -- "
fi

for link in $DOCKER_OPTIONAL_LINKS; do
  IFS=":" <<<"$link" read link_container link_alias
  link_string="--link $link_container:${link_alias:=$link_container}"
  if is_container_running $link_container; then
    DOCKER_LINKS+=" $link_string"
  else
    echo "Container '$link_container' not running, skipping $link_string" >&2
  fi
done

# Attempt to rm existing container with same name. Will fail if running.
if [ -n "$delete_old_container" ]; then
  docker inspect $CONTAINER_NAME >/dev/null 2>&1 && docker rm $CONTAINER_NAME
fi

##
# Invoke `docker run`.

set -e
set -x

docker run \
  $DOCKER_OPTS \
  $DOCKER_VOLUMES \
  $DOCKER_LINKS \
  $DOCKER_ENVS \
  $DOCKER_IMAGE \
  $DOCKER_CMD_PREFIX \
  "${DOCKER_CMD[@]}"
