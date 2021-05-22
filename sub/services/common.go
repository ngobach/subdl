package services

import "strings"

var subtitleFileExt = []string{
	".srt",
}

func filterSubtitleFiles(fileNames []string) []string {
	result := []string{}
	for _, filename := range fileNames {
		for _, ext := range subtitleFileExt {
			if strings.HasSuffix(filename, ext) {
				result = append(result, filename)
				break
			}
		}
	}

	return result
}
