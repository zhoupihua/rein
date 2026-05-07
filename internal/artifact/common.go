package artifact

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type TaskID struct {
	Phase int
	Seq   int
}

func ParseTaskID(s string) (TaskID, bool) {
	var phase, seq int
	n, err := fmt.Sscanf(s, "%d.%d", &phase, &seq)
	if err != nil || n != 2 {
		return TaskID{}, false
	}
	return TaskID{Phase: phase, Seq: seq}, true
}

func (id TaskID) String() string {
	return fmt.Sprintf("%d.%d", id.Phase, id.Seq)
}

func (id TaskID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

type SubTaskID struct {
	Parent TaskID
	Index  int
}

func ParseSubTaskID(s string) (SubTaskID, bool) {
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return SubTaskID{}, false
	}
	parent, ok := ParseTaskID(parts[0] + "." + parts[1])
	if !ok {
		return SubTaskID{}, false
	}
	var idx int
	n, err := fmt.Sscanf(parts[2], "%d", &idx)
	if err != nil || n != 1 || idx < 0 {
		return SubTaskID{}, false
	}
	return SubTaskID{Parent: parent, Index: idx}, true
}

func (id SubTaskID) String() string {
	return fmt.Sprintf("%s.%d", id.Parent.String(), id.Index)
}

func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(string(data), "\r\n", "\n"), nil
}

func WriteFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

type ArtifactFile struct {
	DatePrefix string
	Name       string
	Path       string
}
