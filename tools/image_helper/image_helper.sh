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
        *)
            util::error "Unexpected command '$cmd'."
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

    local img_tar="$1"
    local sbom_path="$2"
    shift 2
    repo_tags_to_push=()
    user_provided_repo_tags=("$@")
    mapfile -t sbom_repo_tags < \
        <("$JQ" -r '.source.target.tags[]' "$sbom_path")

    if [ "${#user_provided_repo_tags[@]}" -lt 1 ]; then
        # push repo tags from SBOM
        repo_tags_to_push=("${sbom_repo_tags[@]}")
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
                    repo_tags_to_push+=("${sbom_repo}${uprt}")
                done
            elif [[ "$uprt" =~ $user_provided_repo_and_tag_regex ]]; then
                # user-provided repo and tag
                repo_tags_to_push+=("$uprt")
            elif [[ "${uprt: -1}" == ":" ]]; then
                # user-provided repo, SBOM tags
                for sbom_tag in "${sbom_tags[@]}"; do
                    repo_tags_to_push+=("${uprt}${sbom_tag}")
                done
            elif [[ "$uprt" =~ $user_provided_registry_regex ]]; then
                # user-provided registry, SBOM path and tags
                for sbom_path in "${sbom_paths[@]}"; do
                    for sbom_tag in "${sbom_tags[@]}"; do
                        repo_tags_to_push+=("${uprt}/${sbom_path}:${sbom_tag}")
                    done
                done
            else
                util::error "Could not determine tagging mode for '$uprt'."
                exit 1
            fi
        done
    fi

    crane_cmd=("$CRANE" "push")
    if [ -n "${CRANE_FLAGS:-}" ]; then
        CRANE_FLAGS=("${CRANE_FLAGS}")
        if [ "${#CRANE_FLAGS[@]}" -gt 0 ]; then
            crane_cmd+=("${CRANE_FLAGS[@]}")
        fi
    fi

    util::info "Pushing '$img_tar' as ${repo_tags_to_push[*]}"
    for rttp in "${repo_tags_to_push[@]}"; do
        "${crane_cmd[@]}" "$img_tar" "$rttp"
        util::success "Pushed '$img_tar' as '$rttp'"
    done
    util::success "Pushed all tags for '$img_tar'"
}

# define utils
util::info() {
    printf "💡 %s\n" "$@"
}

util::warn() {
    printf "⚠️ %s\n" "$@"
}

util::error() {
    printf "❌ %s\n" "$@"
}

util::success() {
  printf "✅ %s\n" "$@"
}


# exec main
main "$@"
