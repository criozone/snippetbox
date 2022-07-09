{{ template "base" . }}

{{ define "title" }}Create a New Snippet{{end}}

{{define "main"}}
<form action="/snippet/create" method="POST" >
    <div>
        <label for="title">Title:</label>
        <input type="text" name="title" id="title" />
    </div>
    <div>
        <label for="content">Content:</label>
        <textarea id="content" name="content"></textarea>
    </div>
    <div>
        <label>Delete in:
            <input type="radio" name="expires" value="365" checked> One year
            <input type="radio" name="expires" value="7" checked> One yweek
            <input type="radio" name="expires" value="1" checked> One day
        </label>
    </div>
    <div>
        <input type="submit" value="Publish snippet">
    </div>

</form>
{{end}}