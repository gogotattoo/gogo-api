<html>
<head>
       <title>Upload file</title>
</head>
<body>
<style type="text/css">
.lab{width:10em;text-align:right;display:inline-block;margin-right:1em;}
.fl{list-style:none}
</style>
<form enctype="multipart/form-data" action="/upload" method="post">
 <ul class="fl">
  <li>
    <label for="artist_name">Who made it?</label>
    <input class="lab" type="text" name="artist_name" value="gogo"/>
  </li>
  <li>
   <label for="made_at">Where was it made?</label>
   <input class="lab" type="text" name="made_at" value="chushangfeng"/>
  </li>
  <li>
    <label for="made_date">When was it made?</label>
    <input class="lab" type="text" name="made_date" value="{{ .Date }}"/>
  </li>

    <input type="file" name="uploadfile" />
    <input type="hidden" name="token" value="{{ .Token }}"/>
    <input type="submit" value="upload" />
  </ul>
</form>
</body>
</html>
