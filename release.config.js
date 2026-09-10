module.exports = {
    branches: [{name: 'main'}],
    verifyConditions: [
        "@semantic-release/github"
    ],
    prepare: [],
    publish: [
        "@semantic-release/github"
    ],
};
