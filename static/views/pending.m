<div class="main">
  <h1>Pending</h1>
  <ul class="results">
    {{#results}}
      <li>
	<div class="link">
          <a href="{{Url}}">{{Title}}</a> 
            <form action="/_pending" method="post" style="display:inline;">
              <input type="hidden" name="id" value="{{Id}}"/>
              <input type="hidden" name="url" value="{{Url}}"/>
              <button value="submit">Approve</button>
            </form>
        </div>
        <div class="domain">{{Domain}}</div>
      </li>
    {{/results}}
  </ul>
</div>