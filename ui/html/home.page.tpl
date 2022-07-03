{{template "base" .}}

{{define "title"}}Home{{end}}

{{define "main"}}
    <h2>Latest Snippets</h2>
    {{if .Snippets}}
    <table>
        <tr>
            <th>Title</th>
            <th>Created</th>
            <th>ID</th>
        </tr>
        {{range .Snippets}}
        <tr>
            <td><a href="/snippet?id={{.Id}}">{{.Title}}</a></td>
            <td>{{humanDate .Created}}</td>
            <td>#{{.Id}}</td>
        </tr>
        {{end}}
    </table>
    {{else}}
    <p>There`s nothing to see here... yet!</p>
    {{end}}
{{end}}
