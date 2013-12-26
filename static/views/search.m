<div class="main">
  <div class="search">
    <form action="/search" method="get">
      <button value="submit">Search</button>
      <div><input class="input" type="text" name="q" id="q" value="{{query}}"/></div>
    </form>
  </div>
  <h5>{{total}} results</h5>
  {{#info}}
    <h3>Problem while searching {{.}}</h3>
  {{/info}}
  <ul class="results">
    {{#results}}
      <li>{{> _script.m}}</li>
    {{/results}}
  </ul>
</div>
