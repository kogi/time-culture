const displayBtn = document.getElementById("display");
const opeBtn = document.getElementById("ope");
const controlBtn = document.getElementById("control");

displayBtn.addEventListener("click", function(e) {
    location.href = "/display";
})

opeBtn.addEventListener("click", function(e) {
    location.href = "/operation";
})

controlBtn.addEventListener("click", function(e) {
    location.href = "/control";
})