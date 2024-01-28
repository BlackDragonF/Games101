package loader

import (
	"assignment1/common"
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Load(path string) []string {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		fmt.Sprintln(os.Stderr, err)
	}
	scanner := bufio.NewScanner(f)
	lines := make([]string, 0)
	// Load file
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Sprintln(os.Stderr, err)
	}
	return lines
}

// Support vertex and face only
func ParseCount(lines []string) (int, int, int, int) {
	vNum, fNum, vtNum, vnNum := 0, 0, 0, 0
	for _, line := range lines {
		words := strings.Split(line, " ")
		switch words[0] {
		case "f":
			fNum++
		case "v":
			vNum++
		case "vt":
			vtNum++
		case "vn":
			vnNum++
		}
	}
	return vNum, fNum, vtNum, vnNum
}

func ParseVertex(vNum int, lines []string) ([]string, []common.Vec3f, error) {
	vertexs := make([]common.Vec3f, 0)
	newLines := make([]string, 0)

	for _, line := range lines {
		words := strings.Split(line, " ")
		switch words[0] {
		case "v":
			x, err := strconv.ParseFloat(words[1], 64)
			y, err := strconv.ParseFloat(words[2], 64)
			z, err := strconv.ParseFloat(words[3], 64)
			if err != nil {
				return nil, nil, err
			}
			var w float64 = 1
			if len(words) == 5 {
				w, err = strconv.ParseFloat(words[4], 64)
				if err != nil {
					return nil, nil, err
				}
				if w == 0 {
					return nil, nil, errors.New("w cannot be 0")
				}
				x /= w
				y /= w
				z /= w
			}
			vertexs = append(vertexs, common.Vec3f{x, y, z})
		default:
			newLines = append(newLines, line)
		}
	}
	if len(vertexs) != vNum {
		return nil, nil, errors.New("wrong number of vertexs")
	}
	return newLines, vertexs, nil
}

func ParseTexture(vtNum int, lines []string) ([]string, []common.Vec3f, error) {
	textures := make([]common.Vec3f, 0)
	newLines := make([]string, 0)

	for _, line := range lines {
		words := strings.Split(line, " ")
		switch words[0] {
		case "vt":
			if len(words) > 4 {
				return nil, nil, errors.New("wrong number of texture")
			}
			u, v, w := 0., 0., 0.
			u, err := strconv.ParseFloat(words[1], 64)
			if len(words) == 3 {
				v, err = strconv.ParseFloat(words[2], 64)
			}
			if len(words) == 4 {
				w, err = strconv.ParseFloat(words[3], 64)
			}
			if err != nil {
				return nil, nil, err
			}
			textures = append(textures, common.Vec3f{u, v, w})
		default:
			newLines = append(newLines, line)
		}
	}
	if len(textures) != vtNum {
		return nil, nil, errors.New("wrong number of texture")
	}
	return newLines, textures, nil
}

func ParseNormal(vnNum int, lines []string) ([]string, []common.Vec3f, error) {
	normals := make([]common.Vec3f, 0)
	newLines := make([]string, 0)

	for _, line := range lines {
		words := strings.Split(line, " ")
		switch words[0] {
		case "vn":
			if len(words) != 4 {
				return nil, nil, errors.New("wrong number of normals")
			}
			x, err := strconv.ParseFloat(words[1], 64)
			y, err := strconv.ParseFloat(words[2], 64)
			z, err := strconv.ParseFloat(words[3], 64)
			if err != nil {
				return nil, nil, err
			}
			sum := x + y + z
			normals = append(normals, common.Vec3f{x / sum, y / sum, z / sum})
		default:
			newLines = append(newLines, line)
		}
	}
	if len(normals) != vnNum {
		return nil, nil, errors.New("wrong number of normalss")
	}
	return newLines, normals, nil
}

func ParseFace(fNum int, lines []string) ([]string, []common.Vec3i, error) {
	faces := make([]common.Vec3i, 0)
	newLines := make([]string, 0)

	for _, line := range lines {
		words := strings.Split(line, " ")
		switch words[0] {
		case "f":
			if len(words) > 4 {
				return nil, nil, errors.New("only support triangle face")
			}
			ele1 := strings.Split(words[1], "/")
			ele2 := strings.Split(words[2], "/")
			ele3 := strings.Split(words[3], "/")
			switch len(ele1) {
			case 1:
				ind1, err := strconv.ParseInt(ele1[0], 10, 64)
				ind2, err := strconv.ParseInt(ele2[0], 10, 64)
				ind3, err := strconv.ParseInt(ele3[0], 10, 64)
				if err != nil {
					return nil, nil, err
				}
				faces = append(faces, common.Vec3i{ind1 - 1, ind2 - 1, ind3 - 1})
			case 2:
				fmt.Sprintln(os.Stderr, "Warning: not support vt for now.")
				ind1, err := strconv.ParseInt(ele1[0], 10, 64)
				ind2, err := strconv.ParseInt(ele2[0], 10, 64)
				ind3, err := strconv.ParseInt(ele3[0], 10, 64)
				if err != nil {
					return nil, nil, err
				}
				faces = append(faces, common.Vec3i{ind1 - 1, ind2 - 1, ind3 - 1})
			case 3:
				fmt.Sprintln(os.Stderr, "Warning: not support vt and vn for now.")
				ind1, err := strconv.ParseInt(ele1[0], 10, 64)
				ind2, err := strconv.ParseInt(ele2[0], 10, 64)
				ind3, err := strconv.ParseInt(ele3[0], 10, 64)
				if err != nil {
					return nil, nil, err
				}
				faces = append(faces, common.Vec3i{ind1 - 1, ind2 - 1, ind3 - 1})
			}
		default:
			newLines = append(newLines, line)
		}
	}
	if len(faces) != fNum {
		return nil, nil, errors.New("wrong number of faces")
	}
	return newLines, faces, nil
}
