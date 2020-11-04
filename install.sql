create database Calculator2; use Calculator2;
create table Users( PrimaryKey int NOT NULL AUTO_INCREMENT, Username varchar( 100 ), Password varchar( 100 ), PRIMARY KEY( PrimaryKey )  );
create table Calculations( PrimaryKey int NOT NULL AUTO_INCREMENT, X varchar( 20 ), Op char(1), Y varchar(20), Val varchar(20), UserKey int, Date datetime, PRIMARY KEY( PrimaryKey ) );