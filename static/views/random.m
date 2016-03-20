<center class="center">
  <p>
  {{#results}}
  <div class="row">
    <div class="col-md-8 col-md-offset-2">
      <h1><a href="/url?s={{Short}}&url={{Uri}}">{{Title}}</a></h1>
      <div class="domain">{{Domain}}</div>
      <div>
        <a href="/s/{{Short}}" onclick="shareLink(this); return false;" style="font-size:0.8em">share</a>
      </div>
      <div>
        <input type="text" value="foo" style="display: none;">
      </div>
    </div>
  </div>
  {{/results}}
</center>
