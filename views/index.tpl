<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{.Title}}</title>
</head>
<body>
    <h1>{{.Content}}</h1>

    {{if .User}}
        <p>Welcome, {{.User}}!</p>
    {{else}}
        <p>Welcome, Guest!</p>
    {{end}}
    <p>Upload file to remove background</p>
    <form enctype="multipart/form-data" action="/upload" method="POST">
        <input type="file" name="file" accept="image/*" required>
        <button type="submit">Upload</button>
    </form>
    {{if .Error}}
    <p style="color: red;">{{.Error}}</p>
    {{else}}
        <p>Uploaded successfully, Waiting for processing</p>
    {{end}}
    
</body>
</html>
