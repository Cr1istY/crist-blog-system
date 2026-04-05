package utils

import (
	"errors"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mozillazg/go-pinyin"
)

var (
	nonAlnumRe    = regexp.MustCompile(`[^a-z0-9\s-]`)
	spaceOrDashRe = regexp.MustCompile(`[\s-]+`)
	chineseCharRe = regexp.MustCompile(`[\x{4e00}-\x{9fff}]`)
)

func ExtractPostTitle(content string) string {
	index := strings.Index(content, " ")
	if index == -1 {
		return ""
	}
	afterHash := content[index:]
	lines := strings.Split(afterHash, "\n")
	title := strings.TrimSpace(lines[0])
	return title
}

func ToSlug(title string) (string, error) {
	if title == "" {
		return "", errors.New("title is required when generating slug")
	}
	slug := toSlugWithoutRandom(title)
	if slug == "" {
		return slug, errors.New("slug is empty, please check your title")
	}
	return slug, nil
}

func toSlugWithoutRandom(title string) string {
	if title == "" {
		return "untitled"
	}
	var sb strings.Builder
	for _, c := range title {
		// 如果是中文，转为拼音
		if chineseCharRe.MatchString(string(c)) {
			args := pinyin.NewArgs()
			args.Style = pinyin.Normal
			py := pinyin.Pinyin(string(c), args)
			// 多音字，取第一个读音
			if len(py) > 0 && len(py[0]) > 0 {
				sb.WriteString(py[0][0])
			}
		} else {
			// 非中文字符，直接保留
			sb.WriteRune(c)
		}
	}
	slug := strings.ToLower(sb.String())
	slug = nonAlnumRe.ReplaceAllString(slug, "")
	slug = spaceOrDashRe.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = "untitled"
	}
	return slug
}

func SlugToSlugWithRandom(slug string) string {
	// 生成随机数，填充在slug之后，防止slug重复
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Intn(25565)
	numChar := strconv.Itoa(num)
	slug += "-"
	slug += numChar
	return slug
}
