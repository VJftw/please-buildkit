#!/usr/bin/env bash
#
set -Eeuo pipefail

main() {
    cmd="$1"
    shift 1

    case "$cmd" in
        push)
            push "$@"
        ;;
        update_refs_in_file)
            update_refs_in_file "$@"
        ;;
        *)
            log::error "Unexpected command '$cmd'."
            exit 1
        ;;
    esac
}

push() {
    # Support pushing by user-provided repository and tag:
    #   `push index.docker.io/foo/bar:my-tag`
    # Support pushing by user-provided repository, SBOM repository tags:
    #   `push index.docker.io/foo/bar:`
    # Support pushing by SBOM repository, user-provided tag:
    #   `push :my-tag`
    # Support pushing by user-provided registry, SBOM repository path and tags:
    #   `push localhost:5000`

    img_tar="$(util::parse_flag img_tar "$@")"
    sbom_path="$(util::parse_flag sbom_path "$@")"
    IFS=',' read -r -a user_provided_repo_tags <<< "$(util::parse_flag user_provided_repo_tags "$@")"

    repo_tags_to_push=($(get_repo_tags "$sbom_path" "${user_provided_repo_tags[@]}"))

    crane_cmd=("$CRANE" "push")
    if [ -n "${CRANE_FLAGS:-}" ]; then
        CRANE_FLAGS=("${CRANE_FLAGS}")
        if [ "${#CRANE_FLAGS[@]}" -gt 0 ]; then
            crane_cmd+=("${CRANE_FLAGS[@]}")
        fi
    fi

    log::info "Pushing '$img_tar' as ${repo_tags_to_push[*]}"
    for rttp in "${repo_tags_to_push[@]}"; do
        "${crane_cmd[@]}" "$img_tar" "$rttp"
        log::success "Pushed '$img_tar' as '$rttp'"
    done
    log::success "Pushed all tags for '$img_tar'"
}

update_refs_in_file() {
    file="$(util::parse_flag file "$@")"
    sbom_path="$(util::parse_flag sbom_path "$@")"
    IFS=',' read -r -a aliases <<< "$(util::parse_flag aliases "$@")"
    IFS=',' read -r -a user_provided_repo_tags <<< "$(util::parse_flag user_provided_repo_tag "$@")"
    log::info "Updating image refs in file '$file' for (${aliases[*]})"

    repo_tags=($(get_repo_tags "$sbom_path" "${user_provided_repo_tags[@]}"))
    # prioritise srcdigest- tag, otherwise use any other tag.
    repo_tag="${repo_tags[0]}"
    if printf "%s\n" "${repo_tags[@]}" | grep "srcdigest-" > /dev/null; then
        repo_tag="$(printf "%s\n" "${repo_tags[@]}" | grep "srcdigest-" | head -n1)"
    fi

    # loop through the aliases and replace in file
    for alias in "${aliases[@]}"; do
        log::info "Replacing '$alias' with '$repo_tag'"
        sed -i "s#${alias}[^\"]*#${repo_tag}#g" "$file"
    done
    log::success "Finished updating image refs in file '$file'"
}

get_repo_tags() {
    local sbom_path="$1"
    shift 1
    local user_provided_repo_tags=("$@")

    repo_tags=()
    user_provided_repo_tags=("$@")
    mapfile -t sbom_repo_tags < \
        <("$JQ" -r '.source.target.tags[]' "$sbom_path")

    if [ "${#user_provided_repo_tags[@]}" -lt 1 ]; then
        # push repo tags from SBOM
        repo_tags=("${sbom_repo_tags[@]}")
    else
        # complete user provided repo tags
        mapfile -t sbom_repos < \
            <(printf "%s\n" "${sbom_repo_tags[@]}" | cut -f1 -d: | sort -u)
        mapfile -t sbom_tags < \
            <(printf "%s\n" "${sbom_repo_tags[@]}" | cut -f2 -d: | sort -u)
        mapfile -t sbom_paths < \
            <(printf "%s\n" "${sbom_repo_tags[@]}" | cut -f1 -d: | cut -f2- -d/ | sort -u)

        user_provided_repo_and_tag_regex="^.+:?.+\/.+:.+"
        user_provided_registry_regex="^[^:/]+:?[^:]*$"

        for uprt in "${user_provided_repo_tags[@]}"; do
            if [ "${uprt:0:1}" == ":" ]; then
                # SBOM repository, user-provided tag
                for sbom_repo in "${sbom_repos[@]}"; do
                    repo_tags+=("${sbom_repo}${uprt}")
                done
            elif [[ "$uprt" =~ $user_provided_repo_and_tag_regex ]]; then
                # user-provided repo and tag
                repo_tags+=("$uprt")
            elif [[ "${uprt: -1}" == ":" ]]; then
                # user-provided repo, SBOM tags
                for sbom_tag in "${sbom_tags[@]}"; do
                    repo_tags+=("${uprt}${sbom_tag}")
                done
            elif [[ "$uprt" =~ $user_provided_registry_regex ]]; then
                # user-provided registry, SBOM path and tags
                for sbom_path in "${sbom_paths[@]}"; do
                    for sbom_tag in "${sbom_tags[@]}"; do
                        repo_tags+=("${uprt}/${sbom_path}:${sbom_tag}")
                    done
                done
            else
                log::error "Could not determine tagging mode for '$uprt'."
                exit 1
            fi
        done
    fi

    printf "%s\n" "${repo_tags[@]}"
}

# define utils
log::debug() {
    if [ -v DEBUG ] || [ -v PKG ]; then
        >&2 printf ">debug: #${BASH_LINENO[0]}> %s\n" "$1"
    fi
}

log::info() {
    >&2 printf "💡 %s\n" "$@"
}

log::warn() {
    >&2 printf "⚠️ %s\n" "$@"
}

log::error() {
    >&2 printf "❌ %s\n" "$@"
}

log::success() {
    >&2 printf "✅ %s\n" "$@"
}

util::csv_to_array() {
    local csv="$1"

    echo "${csv//,/ }"
}

util::parse_flag() {
    local name="$1"
    shift
    while test $# -gt 0; do
        case "$1" in
            "--${name}="*)
                value="$(echo "$1" | cut -d= -f2-)"
                echo "$value"
                log::debug "parsed flag '${name}': $value (from: \"$*\")"
                shift
            ;;
            *)
                shift
            ;;
        esac
    done
}


# exec main
main "$@"
