const { spawn } = require("child_process");

const activities = ["walking", "running", "cycling", "driving", "transport"];

console.log("What's your favorite activity?")
const gum = spawn("gum", ["choose", ...activities]);

gum.stderr.pipe(process.stderr);

gum.stdout.on("data", data => {
    const activity = data.toString().trim();
    console.log(`I like ${activity} too!`);
});
