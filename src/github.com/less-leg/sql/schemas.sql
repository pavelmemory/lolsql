CREATE TABLE django_admin_log (
   id int(11) NOT NULL AUTO_INCREMENT,
   action_time datetime(6) NOT NULL,
   object_id longtext,
   object_repr varchar(200) NOT NULL,
   action_flag smallint(5) unsigned NOT NULL,
   change_message longtext NOT NULL,
   content_type_id int(11) DEFAULT NULL,
   user_id int(11) NOT NULL,
   PRIMARY KEY (id)
 );