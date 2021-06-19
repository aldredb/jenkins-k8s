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
                    sh 'aws eks --region ap-southeast-1 update-kubeconfig --name cluster-1'                    
                }
            }
        }
        stage('Deploy Updated Image to Cluster'){
            steps {
                sh '''
                    export IMAGE="$registry:$BUILD_NUMBER"
                    sed -ie "s~IMAGE~$IMAGE~g" hello.yaml
                    kubectl apply -f hello.yaml
                    '''
            }
        }
    }
}