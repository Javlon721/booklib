package chatWsHandler

import "hash/fnv"

func hasher(s string) (int, error) {
	h := fnv.New32a()

	_, err := h.Write([]byte(s))

	if err != nil {
		return 0, err
	}

	return int(h.Sum32()), nil
}
