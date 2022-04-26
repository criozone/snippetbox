{{template "base" .}}

{{define "title"}}List{{end}}

{{define "main"}}
    <h2>Snippets list</h2>
    <ul>
        {{range .List}}
            <li>{{.}}</li>
        {{end}}
    </ul>
{{end}}
