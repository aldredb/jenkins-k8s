pipeline {
    agent any
    environment {
        registry = "aldredb/hello"
        registryCredential = 'dockerhub'
        dockerImage = ''
    }
    stages {
         stage('Building image') {
            steps{
                script {
                    dockerImage = docker.build registry + ":$BUILD_NUMBER"
                }
            }
        }
        stage('Upload Image to Docker hub') {
            steps{
                script {
                    docker.withRegistry( '', registryCredential ) {
                        dockerImage.push()
                    }
                }
            }
        }
        stage('Remove Unused docker image') {
            steps{
                sh "docker rmi $registry:$BUILD_NUMBER"
            }
        }
        stage('Update kubeconfig'){
            steps {
                withAWS(region:'ap-southeast-1',credentials:'aws') {
                    sh '''
                       CREDENTIALS=$(aws sts assume-role --role-arn arn:aws:iam::878244765130:role/EksWorkshopCodeBuildKubectlRole --role-session-name jenkins-kubectl --duration-seconds 900)
                       export AWS_ACCESS_KEY_ID="$(echo ${CREDENTIALS} | jq -r '.Credentials.AccessKeyId')"
                       export AWS_SECRET_ACCESS_KEY="$(echo ${CREDENTIALS} | jq -r '.Credentials.SecretAccessKey')"
                       export AWS_SESSION_TOKEN="$(echo ${CREDENTIALS} | jq -r '.Credentials.SessionToken')"
                       export AWS_EXPIRATION=$(echo ${CREDENTIALS} | jq -r '.Credentials.Expiration')
                       aws eks --region ap-southeast-1 update-kubeconfig --name cluster-1
                       '''                    
                }
            }
        }
        stage('Deploy Updated Image to Cluster'){
            steps {
                sh '''
                    export IMAGE="$registry:$BUILD_NUMBER"
                    sed -ie "s~IMAGE~$IMAGE~g" hello.yaml
                    export KUBECONFIG=/var/lib/jenkins/.kube/config
                    kubectl apply -f hello.yaml
                    '''
            }
        }
    }
}