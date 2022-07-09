{{ template "base" . }}

{{ define "title" }}Create a New Snippet{{end}}

{{define "main"}}
<form action="/snippet/create" method="POST" >
    <div>
        <label for="title">Title:</label>
        {{with .Form.Errors.Get "title"}}
        <label class="error">{{.}}</label>
        {{end}}
        <input type="text" name="title" id="title" value='{{ .Form.Get "title" }}' />
    </div>
    <div>
        <label for="content">Content:</label>
        {{with .Form.Errors.Get "content"}}
        <label class="error">{{.}}</label>
        {{end}}
        <textarea id="content" name="content">{{.Form.Get "content"}}</textarea>
    </div>
    <div>
        <label>Delete in:</label>
        {{with .Form.Errors.Get "expires"}}
        <label class="error">{{.}}</label>
        {{end}}
        {{$exp := or (.Form.Get "expires") "365"}}
        <label><input type="radio" name="expires" value="365" {{if (eq $exp "365")}}checked{{end}}> One year</label>
        <label><input type="radio" name="expires" value="7" {{if (eq $exp "7")}}checked{{end}}> One yweek</label>
        <label><input type="radio" name="expires" value="1" {{if (eq $exp "1")}}checked{{end}}> One day</label>
    </div>
    <div>
        <input type="submit" value="Publish snippet">
    </div>

</form>
{{end}}