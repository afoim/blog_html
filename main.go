package main

import (
    "bufio"
    "fmt"
    "github.com/gomarkdown/markdown"
    "github.com/gomarkdown/markdown/html"
    "github.com/gomarkdown/markdown/parser"
    "html/template"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "strings"
    "time"
)

// ShortLink 结构体用于存储短链接信息
type ShortLink struct {
    Code    string // 短链接代码
    PostID  string // 对应的文章ID
}

// Post 结构体用于存储博文信息
type Post struct {
    Title     string
    Content   template.HTML
    Date      time.Time
    Filename  string
    ShortLink string // 短链接代码
}

// PageData 结构体用于存储页面数据
type PageData struct {
    Posts []Post
    Today time.Time
    Post  Post // 用于单篇文章页面
}

// 存储所有短链接
var shortLinks = make(map[string]string)

// 自定义函数将标题转换为锚点 ID
func anchorize(title string) string {
    return strings.ToLower(strings.ReplaceAll(title, " ", "-"))
}

func main() {
    // 创建 dist 目录
    err := os.MkdirAll("dist", 0755)
    if err != nil {
        log.Fatal("创建 dist 目录失败:", err)
    }

    // 读取所有博文
    posts, err := loadPosts()
    if err != nil {
        log.Fatal("加载博文失败:", err)
    }

    // 生成所有页面
    err = generatePages(posts)
    if err != nil {
        log.Fatal("生成页面失败:", err)
    }

    // 生成短链接重定向页面
    err = generateShortLinkPages()
    if err != nil {
        log.Fatal("生成短链接页面失败:", err)
    }

    fmt.Println("网站生成完成！文件保存在 dist 目录中")
}

// 检查第一行是否包含短链接
func parseShortLink(content string) (string, string) {
    scanner := bufio.NewScanner(strings.NewReader(content))
    if scanner.Scan() {
        firstLine := scanner.Text()
        if strings.HasPrefix(firstLine, "!url:") {
            code := strings.TrimSpace(strings.TrimPrefix(firstLine, "!url:"))
            remainingContent := strings.TrimPrefix(content, firstLine+"\n")
            return code, remainingContent
        }
    }
    return "", content
}

// 加载所有博文
func loadPosts() ([]Post, error) {
    var posts []Post

    files, err := ioutil.ReadDir("posts")
    if err != nil {
        return nil, err
    }

    for _, file := range files {
        if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
            content, err := ioutil.ReadFile(filepath.Join("posts", file.Name()))
            if err != nil {
                return nil, err
            }

            // 检查并解析短链接
            shortLinkCode, remainingContent := parseShortLink(string(content))

            // 创建 markdown 解析器
            extensions := parser.CommonExtensions | parser.AutoHeadingIDs
            p := parser.NewWithExtensions(extensions)

            // 创建 HTML 渲染器
            opts := html.RendererOptions{
                Flags: html.CommonFlags | html.HrefTargetBlank,
            }
            renderer := html.NewRenderer(opts)

            // 转换 Markdown 为 HTML
            htmlContent := markdown.ToHTML([]byte(remainingContent), p, renderer)

            // 使用文件名作为标题（去掉.md后缀）
            title := strings.TrimSuffix(file.Name(), ".md")

            // 将 HTML 内容转换为字符串
            htmlStr := string(htmlContent)

            // 替换标题以包含 ID
            for _, heading := range []string{"h1", "h2", "h3", "h4", "h5", "h6"} {
                htmlStr = strings.ReplaceAll(htmlStr, 
                    fmt.Sprintf("<%s>", heading), 
                    fmt.Sprintf("<%s id=\"%s\">", heading, anchorize(title)))
            }

            post := Post{
                Title:     title,
                Content:   template.HTML(htmlStr),
                Date:      file.ModTime(),
                Filename:  title,
                ShortLink: shortLinkCode,
            }

            // 如果存在短链接，添加到映射中
            if shortLinkCode != "" {
                shortLinks[shortLinkCode] = title
            }

            posts = append(posts, post)
        }
    }

    return posts, nil
}

// 生成短链接重定向页面
func generateShortLinkPages() error {
    // 短链接重定向模板
    redirectTemplate := `<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="refresh" content="0;url=/{{.PostID}}.html">
</head>
<body>
    正在重定向到文章页面...
</body>
</html>`

    tmpl, err := template.New("redirect").Parse(redirectTemplate)
    if err != nil {
        return err
    }

    for code, postID := range shortLinks {
        // 创建短链接目录
        err := os.MkdirAll(filepath.Join("dist", code), 0755)
        if err != nil {
            return err
        }

        // 创建 index.html 文件
        file, err := os.Create(filepath.Join("dist", code, "index.html"))
        if err != nil {
            return err
        }

        err = tmpl.Execute(file, struct{ PostID string }{PostID: postID})
        file.Close()
        if err != nil {
            return err
        }
    }

    return nil
}

// 生成所有页面
func generatePages(posts []Post) error {
    // 读取模板文件
    tmpl, err := template.ParseFiles("templates.html")
    if err != nil {
        return err
    }

    // 生成首页
    indexFile, err := os.Create("dist/index.html")
    if err != nil {
        return err
    }
    defer indexFile.Close()

    err = tmpl.Execute(indexFile, PageData{
        Posts: posts,
        Today: time.Now(),
    })
    if err != nil {
        return err
    }

    // 生成每篇文章的页面
    for _, post := range posts {
        file, err := os.Create(filepath.Join("dist", post.Filename+".html"))
        if err != nil {
            return err
        }

        err = tmpl.Execute(file, PageData{
            Post:  post,
            Today: time.Now(),
        })
        file.Close()
        if err != nil {
            return err
        }
    }

    return nil
}