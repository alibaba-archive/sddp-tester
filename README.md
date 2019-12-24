####版本内容  
1.构造了一些敏感数据并自动上传到OSS、RDS、ODPS中,触发SDDP产品的敏感数据检测  
2.模拟黑客的一些异常行为,触发SDDP产品的异常行为检测  
3.可以使用命令一键清理上传用以测试的文件和表
  
####使用说明  
1.测试敏感数据检测  

RDS:  
- step 1.在SDDP的授权页面将测试的RDS库进行授权  
- step 2.检查本机IP是否属于测试RDS的连接白名单  
- step 3.在config.ini中输入RDS的访问账号和密码后保存并关闭  
- step 4.打开命令行界面cd到项目目录  
- step 5.输入./sddp_tester scan -p rds -t <rds类型(mysql/mssql)>  -d <测试数据库名称> -P <数据库端口> -H <数据库域名>   
  
OSS:
- step 1.在SDDP的授权页面将测试的Bucket进行授权  
- step 2.在网页https://usercenter.console.aliyun.com/#/manage/ak中获取账号的AccessKeyId、AccessKeySecret  
- step 3.在config.ini中输入AccessKeyId、AccessKeySecret后保存并关闭  
- step 4.打开命令行界面cd到项目目录  
- step 5.进入到OSS产品中测试Bucket的首页,拿到EndPoint  
- step 6.输入./sddp_tester scan -p oss  --ossendpoint <oss的endpoint> --bucketname <测试的bucket> 
  
ODPS:
- step 1.在SDDP的授权页面将测试的Project进行授权  
- step 2.在网页https://usercenter.console.aliyun.com/#/manage/ak中获取账号的AccessKeyId、AccessKeySecret  
- step 3.在config.ini中输入AccessKeyId、AccessKeySecret后保存并关闭  
- step 4.打开命令行界面cd到项目目录  
- step 5.输入./sddp_tester scan -p odps --odpsendpoint http://service.cn.maxcompute.aliyun.com/api --project <测试的project>
  
2.异常行为模型触发
  
测试指定的实例最好是使用时间较长的,否则一些模型可能会因为没有历史行为基线导致无法产出异常。

RDS:  
- step 1.在SDDP的授权页面将测试的RDS库进行授权  
- step 2.检查本机IP是否属于测试RDS的连接白名单  
- step 3.在config.ini中输入RDS的访问账号和密码后保存并关闭  
- step 4.打开命令行界面cd到项目目录  
- step 5.输入./sddp_tester abnormal -p rds -t <rds类型(mysql/mssql)>  -d <测试数据库名称> -P <数据库端口> -H <数据库域名> --table <敏感表名> --proxy <http://1.1.1.1:8888>

OSS:  
- step 1.在SDDP的授权页面将测试的Bucket进行授权  
- step 2.在网页https://usercenter.console.aliyun.com/#/manage/ak中获取账号的AccessKeyId、AccessKeySecret  
- step 3.在config.ini中输入AccessKeyId、AccessKeySecret后保存并关闭  
- step 4.打开命令行界面cd到项目目录  
- step 5.进入到OSS产品中测试Bucket的首页,拿到EndPoint  
- step 6.输入./sddp_tester abnormal -p oss --ossendpoint <oss的endpoint> --bucketname <测试的bucket> --sensitiveobject <test/1.jpg> --proxyip <http://1.1.1.1:8888>  

ODPS:  
- step 1.在SDDP的授权页面将测试的Project进行授权  
- step 2.在网页https://usercenter.console.aliyun.com/#/manage/ak中获取账号的AccessKeyId、AccessKeySecret  
- step 3.在config.ini中输入AccessKeyId、AccessKeySecret后保存并关闭  
- step 4.打开命令行界面cd到项目目录  
- step 5.输入./sddp_tester abnormal -p odps --odpsendpoint http://service.cn.maxcompute.aliyun.com/api --project <测试的project> --table <敏感表名> --proxy <http://1.1.1.1:8888>
  
3.清理测试文件和表
- step 1.打开命令行界面cd到项目目录  
- step 2.输入./sddp_tester clean
*命令会清理掉上次清理到本次清理之间运行工具产生的所有文件和表 
  
####特别说明  
1.本程序会在您的OSS、ODPS、RDS上传一些用以测试的数据文件,不过在测试结束后,可以通过./sddp_tester clean命令清理数据文件  
2.为了防止文件覆盖,程序可能会产生微许的数据枚举   
3.异常行为模型会下载大量测试bucket的文件  
4.异常行为模型会将测试的Bucket设为公开可访问
5.异常行为模型会将测试的ODPS项目设为未设置标签保护和未设置项目保护