<html>
  <head>
    <meta name="viewport" content="initial-scale=1">
    <meta name="description" content="A screenplay index">
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <title>Got Scripts?</title>
    <style type="text/css">
      html {
        font-family: arial,sans-serif;
	font-size: small;
        color: #404040;
      }
      body {
	padding: 10px 10% 10px 10%;
	min-width: 280px;
        max-width: 960px;
	margin: 0 auto;
      }
      button {
	background: #f1f1f1;
	border: 1px solid #dcdcdc;
        padding:5px;
      }
      ul {
        list-style-type: none;
	padding: 0;
      }
      li {
        display: inline;
	padding-right: 10px;
      }
      ul img {
	vertical-align: bottom;
      }
      .center {
        text-align: center;
	padding-top: 80px;
      }
      .center .search {
        margin: 0 auto;
      }
      .input {
        width: 100%;
        max-width: 600px;
	padding: 5px;
	font-weight: bold;
	border: 1px solid #4285f4;
      }
      .nav {
        padding: 0;
	display: inline-block;
      }
      .nav a {
        text-decoration: none;
	color: #404040;
      }
      .nav a:hover {
        text-decoration: underline;
      }
      .alert {
        width: 100%;
        color: #dd4b39;
	font-weight: bold;
      }
      .results {
        padding: 0;
      }
      .results li {
        display: block;
        padding-bottom: 10px;
      }
      .domain {
        color: #006621;
      }
      .search {
        max-width: 600px;
      }
      .search div {
        overflow: hidden;
        max-width: 550px;
      }
      .search button {
        float: right;
      }
      .search-ico a {
        background: url('/static/image/s.png') no-repeat;
        display: inline-block;
        text-indent:-9999px;
        width: 16px;
        overflow: hidden;
      }
      .home-ico a {
        background: url('/static/image/h.png') no-repeat;
        display: inline-block;
        text-indent:-9999px;
        width: 16px;
        overflow: hidden;
      }
      .search-ico a:hover {
        background: url('/static/image/sh.png') no-repeat;
        cursor: hand;
      }
      .home-ico a:hover {
        background: url('/static/image/hh.png') no-repeat;
        cursor: hand;
      }
      .share {
        display: inline;
        padding-left: 5px;
      }
      .share input {
	padding: 0;
        margin: 0;
        width: 200px;
	border: 1px solid #4285f4;
        background: #f1f1f1;
      }
      .share li {
        display: inline;
      }
    </style>
  </head>
  <body>
    <ul class="nav">
      <li class="home-ico"><a href="/">h</a></li>
      <li><a href="/scripts">latest</a></li>
      <li><a href="/add">+add</a></li>
      <li class="search-ico"><a href="/">s</a></li>
    </ul>
    {{#alert}}<center class="alert">{{alert}}</center>{{/alert}}
    {{{content}}}
    {{> _pager.m}}
    <script>
      function shareLink(obj) {
	href = obj.getAttribute("href");
	input = obj.parentNode.parentNode.getElementsByTagName("input")[0];
	input.style.display = 'inline';
	input.value = location.protocol + '//' + location.hostname + href;
	input.select();
	return false;
      };
    </script>
  </body>
</html>
