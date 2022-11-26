package image

import "strings"

type RepoTag string

func (i RepoTag) GetRepository() string {
	parts := strings.Split(string(i), ":")
	if len(parts) == 2 {
		return parts[0]
	}
	if len(parts) == 3 {
		return strings.Join([]string{parts[0], parts[1]}, ":")
	}

	return ""
}

func (i RepoTag) GetTag() string {
	parts := strings.Split(string(i), ":")
	if len(parts) == 2 {
		return parts[1]
	}
	if len(parts) == 3 {
		return parts[2]
	}

	return ""
}

func (i RepoTag) GetRegistry() string {
	return strings.Split(string(i), "/")[0]
}

func (i RepoTag) GetPath() string {
	return strings.SplitN(i.GetRepository(), "/", 2)[1]
}
