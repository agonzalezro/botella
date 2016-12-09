var languages = [
  "Java",
  "C",
  "C++",
  "Python",
  ".NET",
  "C#",
  "PHP",
  "Javascript",
  "Perl",
  "Objective-c",
  "Ruby",
  "Swift",
  "Visual Basic",
  "Delphi",
  "Go",
  "Groovy",
  "Elixir",
  "Scala",
  "Bash",
]

$(function(){
    $(".language").typed({
        strings: languages,
        typeSpeed: 200,
        backSpeed: 50,
        startDelay: 0,
        backDelay: 500,
        loop: true
    });
});
