{{template "base" .}}

{{define "title"}}Home{{end}}

{{define "main"}}
    <h2>Latest Snippets</h2>
    <ul>
        {{range .Snippets}}
            <li><a href="/snippet?id={{.Id}}">{{.Title}}</a></li>
        {{end}}
    </ul>
{{end}}
