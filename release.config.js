module.exports = {
    branches: [
        {name: 'main'},
        {name: `${process.env.BRANCH}`, prerelease: "rc"}
    ],
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
