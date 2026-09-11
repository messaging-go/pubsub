module.exports = {
    branches: [
        'main',
        {name: '!(main)', prerelease: 'rc'},
    ],
    verifyConditions: [
        "@semantic-release/github"
    ],
    prepare: [],
    publish: [
        "@semantic-release/github"
    ],
};
