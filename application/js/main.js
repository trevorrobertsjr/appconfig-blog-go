const htmlDarkKnightOutput =
  '<section><img src="assets/img/dark_knight_cecil.png" alt="image of Cecil as a Dark Knight"></section><section><p>Cecil continued to walk the path of the Dark Knight. The Path of the Paladin was too new and unfamiliar to Cecil for him to accept. Cecil determined he was strong enough to defeat the King of Baron with the Dark Sword. The spirit of the father of Cecil was disappointed, but he ultimately understood.</p></section>';
const htmlPaladinOutput =
  '<section><img src="assets/img/paladin_cecil.png" alt="image of Cecil as a Paladin"></section><section><p>Cecil heeded the words of the light and took his rightful place as the prophesied Paladin. Contrary to his fears, selecting the light grew his power exponentially. He then realized the King of Baron encouraged Cecil to take the dark sword to prevent him from realizing his full potential.</p></section>';

fetch(
  "http://localhost:2772/applications/blogAppConfigGo/environments/prod/configurations/whichSide?flag=allegiance"
)
  .then((response) => response.json())
  .then((data) => {
    if ((data.choice = "paladin")) {
      document.getElementById("displayCecil").innerHTML = htmlPaladinOutput;
    } else {
      document.getElementById("displayCecil").innerHTML = htmlDarkKnightOutput;
    }
    console.log(data);
  });
