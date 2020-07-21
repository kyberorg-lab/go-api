@Library('common-lib@1.4') _
pipeline {
    agent any;
    stages {
        stage('Docker') {
            steps {
                script {
                    def repo = 'kyberorg/go-api';
                    def tags = [];
                    String tag;
                    String dockerfileDir
                    if (env.BRANCH_NAME.equals("master")) {
                        tag = "stable";
                        dockerfileDir = "docker/prod"
                    } else {
                        tag = env.BRANCH_NAME;
                        dockerfileDir = "docker"
                    }
                    tags << tag;

                    dockerBuild(repo: repo, tags: tags, dockerfileDir: dockerfileDir);
                    dockerLogin(creds: 'hub-docker');
                    dockerPush();
                    dockerLogout();
                    dockerClean();
                }
            }
        }
        stage('Deploy') {
            steps {
                script {
                    String hookUrl;
                    switch (env.BRANCH_NAME) {
                        case "master":
                            hookUrl = "https://docker.yatech.eu/api/webhooks/b2bbdb1f-b4d1-48a6-85f1-f5661572f367?tag=stable";
                            break;
                        default:
                            hookUrl = "https://docker.yatech.eu/api/webhooks/b2bbdb1f-b4d1-48a6-85f1-f5661572f367?tag=" + env.BRANCH_NAME;
                            break;
                    }
                    //no hook - no deploy
                    if(hookUrl.equals('')) { return; }
                    deployToSwarm(hookUrl: hookUrl);
                    sleep(10); //pause for application to be started
                }
            }
        }
    }
    post {
        always {
            cleanWs();
        }
    }
}
