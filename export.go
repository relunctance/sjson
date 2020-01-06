package sjson

func Select(json string, paths []string) (string, error) {
	return SelectString(json, paths)
}

func SelectPath(json string, path string) (string, error) {
	return SelectString(json, []string{path})
}

func SelectBytes(json []byte, paths []string) ([]byte, error) {
	return getByBytes(json, paths)
}

func SelectBytesPath(json []byte, path string) ([]byte, error) {
	return getByBytes(json, []string{path})

}

func SelectString(json string, paths []string) (string, error) {
	data, err := SelectBytes([]byte(json), paths)
	return string(data), err
}
