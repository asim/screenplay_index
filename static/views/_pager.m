{{#pager}}
<ul class="pager">
  <li class="previous {{previousState}}">
   <a href="{{previousPage}}">Previous</a>
  </li>
  <li class="next {{nextState}}">
   <a href="{{nextPage}}">Next</a> 
  </li>
</ul>
{{/pager}}
{{^pager}}
<p class="pager">&nbsp;<p>
{{/pager}}
