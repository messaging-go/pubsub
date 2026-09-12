const cfg = process.env.BRANCH == "main" ? [] : [{name: `${process.env.BRANCH}`, prerelease: "rc"}]
module.exports = {
    branches: [{name: "main"}, ...cfg],
    plugins: [
        "@semantic-release/commit-analyzer",
        "@semantic-release/release-notes-generator",
        "@semantic-release/github"
    ],
    verifyConditions: [
        "@semantic-release/github"
    ],
    publish: [
        "@semantic-release/github"
    ]
}
