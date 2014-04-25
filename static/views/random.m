<div class="center">
  <ul class="random-results">
  {{#results}}
    <li>
      <h1><a href="/url?s={{Short}}&url={{Uri}}">{{Title}}</a></h1>
      <div class="domain">{{Domain}}</div>
      <ul class="share">
        <li><a href="/s/{{Short}}" onclick="shareLink(this); return false;" style="font-size:0.8em">share</a></li>
        <li><input type="text" value="foo" style="display: none;"></li>
      </ul>
    </li>
  {{/results}}
  </ul>
</div>
