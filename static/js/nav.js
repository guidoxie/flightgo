$(document).ready(function() {
   $('#nav-button').on('click', function() {
       $('header').toggleClass('nav-open');
   });
});
document.getElementById("nav-button").addEventListener("click", function() {
   document.getElementById("nav-toggle").classList.toggle("active");
});