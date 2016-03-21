<h3>
  <a href="/url?s={{Short}}&url={{Uri}}">{{Title}}</a>
</h3>

<div>{{Domain}}</div>

<p>
  <small>{{Meta}}</small>
</p>

<ul class="list-inline">
  <li>
    <a href="http://twitter.com/share?text={{Title}}&amp;url=http://scridx.com/s/{{Short}}"
        onclick="window.open(this.href, 'twitter-share', 'width=550,height=235');return false;">
        <i class="fa fa-twitter"></i>
        <span class="hidden">Twitter</span>
    </a>
  <li>
  </li>
    <a href="https://www.facebook.com/sharer/sharer.php?u=http://scridx.com/s/{{Short}}"
        onclick="window.open(this.href, 'facebook-share','width=580,height=296');return false;">
        <i class="fa fa-facebook-official"></i>
        <span class="hidden">Facebook</span>
    </a>
  <li>
  </li>
    <a href="https://plus.google.com/share?url=http://scridx.com/s/{{Short}}"
       onclick="window.open(this.href, 'google-plus-share', 'width=490,height=530');return false;">
        <i class="fa fa-google-plus"></i>
        <span class="hidden">Google+</span>
    </a>
  </li>
  <li>
    <a href="/s/{{Short}}" onclick="shareLink(this); return false;" style="font-size:0.8em">
      <i class="fa fa-clipboard"></i>
    </a>
    &nbsp;
    <span><input type="text" value="foo" style="display: none;"></span>
  </li>
</ul>


