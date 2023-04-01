output "eks-spoke1-endpoint" {
    value = module.eks-spoke1.cluster_endpoint
}
output "eks-spoke2-endpoint" {
    value = module.eks-spoke2.cluster_endpoint
}
output "nyancat_repo_url" {
    description = "URL of the nyancat respository"
    value = aws_ecr_repository.nyancat.repository_url
}
output "nyancat_registry_id" {
    description = "URL of the nyancat respository"
    value = aws_ecr_repository.nyancat.registry_id
}
output "nginx_repo_url" {
    description = "URL of the Nginx respository"
    value = aws_ecr_repository.nginx.repository_url
}
output "nginx_registry_id" {
    description = "URL of the Nginx respository"
    value = aws_ecr_repository.nginx.registry_id
}