const navbar = document.querySelector("#navbar")
const menu = document.querySelector("#menubtn")

menu.addEventListener("click", () => {
    if (navbar.classList.contains("-ml-52")) {
        navbar.classList.remove("-ml-52");
        navbar.classList.add("ml-0");
    } else {
        navbar.classList.remove("ml-0");
        navbar.classList.add("-ml-52");
    }
    let closed = navbar.classList.contains('-ml-52');
    if (closed) {
        menu.classList.remove("fa-xmark");
        menu.classList.add("fa-bars");
    } else {
        menu.classList.remove("fa-bars");
        menu.classList.add("fa-xmark");
    }
})

window.addEventListener("touchstart", e => {
  if (!navbar.contains(e.target) && (!menu.contains(e.target))){
  navbar.classList.remove("ml-0");
  navbar.classList.add("-ml-52");
  menu.classList.remove("fa-xmark");
  menu.classList.add("fa-bars");
} 
});

let map;

async function initMap() {
  const { Map } = await google.maps.importLibrary("maps");

  map = new Map(document.getElementById("map"), {
    center: { lat: 53.541190, lng: -113.295629 },
    zoom: 12,      
  });
}

initMap();