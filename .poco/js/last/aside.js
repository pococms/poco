var scrollTimer = -1;
function bodyScroll()
{
  if(scrollTimer != -1){
    clearTimeout(scrollTimer);
  }
  scrollTimer = window.setTimeout(sidebarHeight, 100);
}
function sidebarHeight() {
    console.log("sidebarHeight()")
  s=document.getElementById('aside-poco');
  a=document.getElementById('article-poco'); 
  if (s != null && a != null) {
    // If there's a sidebar,
    ha=a.offsetHeight;
    hs=s.offsetHeight;
    // If the article is longer, make
    // the sidebar as long as the article.
    if (ha>hs) {
      h=a.offsetHeight+'px';s.style.height=h;
    }
    // If the sidebar is longer than the
    // article, make the article as long
    // as the sidebar.
    else{
      h=s.offsetHeight+'px';a.style.height=h;
    }
  }

} 
document.onreadystatechange = function () {
  if (document.readyState == "interactive") {
    // Init or start code here
    bodyScroll();
  }
} 
window.onresize=sidebarHeight;
/* ES5 works fine */
document.getElementsByTagName('body')[0].onscroll=bodyScroll();


