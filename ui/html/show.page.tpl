{{template "base" .}}

{{define "title"}}Snippet #{{.Snippet.Id}}{{end}}

{{define "main"}}
<div class='snippet'>
    <div class='metadata'>
        <strong>{{.Snippet.Title}}</strong>
        <span>#{{.Snippet.Id}}</span>
    </div>
    <pre><code>{{.Snippet.Content}}</code></pre>
    <div class='metadata'>
        <time>Created: {{.Snippet.Created}}</time>
        <time>Expires: {{.Snippet.Expires}}</time>
    </div>
</div>
{{end}}