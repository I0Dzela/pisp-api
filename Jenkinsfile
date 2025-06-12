pipeline {
    agent any
    stages {
        stage('[PISP_SPECS] Docker image build') {
            steps {
                script {
                    withCredentials([
                        sshUserPrivateKey(credentialsId: 'pisp-repository',
                                          keyFileVariable: 'keyFile',
                                          passphraseVariable: 'passphrase',
                                          usernameVariable: 'username')
                    ]) {
                        writeFile file: './docker/secrets/repository-credentials', text: readFile(keyFile)
                        sh '''
                        docker compose build
                        '''
                    }
                }
            }
        }

        stage('[PISP_SPECS] Docker image publish') {
            steps {
                script {
                    withCredentials([usernamePassword(credentialsId: 'local-docker-registry',
                                                      usernameVariable: 'username',
                                                      passwordVariable: 'password')]) {
                        sh '''
                        docker login -u ${username} -p ${password} localhost:5000
                        docker tag pisp/specs:1.0.0 localhost:5000/pisp/specs:1.0.0 && \
                        docker push localhost:5000/pisp/specs:1.0.0
                        '''
                    }
                }
            }
        }

        stage ('[PISP_SPECS] Restart deployment') {
            steps {
                build job: 'pisp-deployment-3', propagate: true, wait: false
            }
        }
    }
}