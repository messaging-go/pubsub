module.exports = {
    branches: [
        {name: 'main'},
        {name: '**', prerelease: 'rc'},
    ],
    verifyConditions: [
        "@semantic-release/github"
    ],
    prepare: [],
    publish: [
        "@semantic-release/github"
    ],
};
