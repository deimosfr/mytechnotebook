---
weight: 999
url: "/CAS_\\:_Mise_en_place_d'un_serveur_SSO/"
title: "CAS: Setting Up an SSO Server"
description: "Learn how to set up and configure Central Authentication Service (CAS) for Single Sign-On (SSO) across multiple applications, with specific instructions for Confluence and Jira integration."
categories: ["Servers", "Security", "Authentication"]
date: "2009-10-08T16:59:00+02:00"
lastmod: "2009-10-08T16:59:00+02:00"
tags: ["cas", "sso", "authentication", "ldap", "tomcat", "confluence", "jira"]
toc: true
---

## Introduction

[Single sign-on (SSO)](https://en.wikipedia.org/wiki/Single_sign_on) is a property of access control of multiple, related, but independent software systems. With this property a user logs in once and gains access to all systems without being prompted to log in again at each of them. Single sign-off is the reverse property whereby a single action of signing out terminates access to multiple software systems.

As different applications and resources support different authentication mechanisms, single sign-on has to internally translate to and store different credentials compared to what is used for initial authentication.

## Purpose

In order to avoid users having to sign on for each application, we want to have them authenticate on a dedicated application. Then, each time the user accesses another application, this application will ask the authentication application if the user is already authenticated and permitted to use this application.

The product chosen for this purpose is [JAS-SIG CAS](https://www.jasig.org/). It's widely used in various environments, from big universities to SMEs or large companies such as Valtech, Smile or CGG Veritas.

The 2 main applications currently used in my company (confluence & jira) have already been "CASified" by other people and the process is well documented. For other applications, there are several client libraries developed. These libraries must be integrated and used in the various applications. The difficulty of this depends on the application's source availability, how clear the sources are, and what knowledge admins and/or developers have of the application. Successful CASification of some applications will be described on this page.

## Installation

### Tomcat installation & configuration

You need a tomcat server, preferably installed on a debian system (but it can really be any kind of system).

So, install a basic debian server environment.

Install sun java 5 or 6 (it doesn't really matter):

```bash
apt-get install sun-java6-jdk sun-java6-fonts
```

Install tomcat (release 5.5 actually).
To avoid problems, you'll need to edit `/etc/default/tomcat5.5` and uncomment the line which defines TOMCAT5_SECURITY. By default, this variable is set to "yes". Set it to "no".

You'll also need to have SSL enabled on your tomcat. For this, follow the instructions given on these pages (just change the validity length from 365 to 3652):
http://blogs.dfwikilabs.org/pigui/2007/12/10/configuring-tomcat-55-for-ssl-using-openssl/
http://tomcat.apache.org/tomcat-5.5-doc/ssl-howto.html

You'll have to use the certificate on the various clients. I'll remind you this when needed.

### JA-SIG Webapp configuration & compilation

You have to download and compile the cas web application (the cas server itself). For this you'll need to use maven. For convenience and clarity, perform the following operations on a development server. You don't have to install and compile on the production server.

Install maven2:

```bash
apt-get install maven2
```

If your system is not debian lenny, just help yourself.

Download the tar file from ja-sig website (http://www.jasig.org/cas/download).

Untar the file where you want.

- Change to the directory you've just untared and edit the file pom.xml

At the end of the dependencies element, add this:

```xml
<dependency>
    <groupId>${project.groupId}</groupId>
    <artifactId>cas-server-support-ldap</artifactId>
    <version>${project.version}</version>
</dependency>
```

- A few lines below, comment out the line

```xml
               <!--
                <module>cas-server-support-ldap</module>
                -->
```

- Save the file and exit your editor
- Then change directory to cas-server-webapp/src/main/webapp/WEB-INF and edit the file deployerConfigContext.xml
- Find the string "SimpleTestUsernamePassword". Comment out this bean and add a new bean for the ldap authentication. You should have this:

```xml
<!--
    <bean class="org.jasig.cas.authentication.handler.support.SimpleTestUsernamePasswordAuthenticationHandler" />
-->
    <bean class="org.jasig.cas.adaptors.ldap.BindLdapAuthenticationHandler">
                <property name="filter" value="uid=%u" />
                <property name="searchBase" value="dc=openldap,dc=mycompany,dc=lan" />
                <property name="contextSource" ref="contextSource" />
    </bean>
```

- In the same file, go to the end and before the &lt;/beans&gt;, add this new bean declaration:

```xml
      <bean id="contextSource" class="org.springframework.ldap.core.support.LdapContextSource">
                <property name="anonymousReadOnly" value="false" />
                <property name="pooled" value="true" />
                <property name="urls">
                        <list>
                                <value>ldap://tasmania/</value>
                                <value>ldap://star1/</value>
                        </list>
                </property>
                <!-- uncomment the following lines if you nead to bind to the directory. Adjust to your needs.
                <property name="password" value="{password_goes_here}" />
                <property name="userName" value="{username_goes_here}" />
                <property name="baseEnvironmentProperties">
                        <map>
                                <entry>
                                        <key><value>java.naming.security.protocol</value></key>
                                        <value>ssl</value>
                                </entry>
                                <entry>
                                        <key><value>java.naming.security.authentication</value></key>
                                        <value>simple</value>
                                </entry>
                        </map>
                </property>
                -->
        </bean>
```

- Save the file and exit your editor
- Go back to the top directory (the one you untared the initial file) and compile the whole stuff with maven:

```bash
mvn -Dmaven.test.skip=true package install
```

If the compilation succeeded, you have a cas.war file in $TOPDIR/cas-server-webapp/target.

### JA-SIG Webapp installation & tests

Copy the war file to the webapps directory (`/var/lib/tomcat5.5/webapps`) of your tomcat server.

Once the file has been copied, eventually restart tomcat. Then, if everything is ok, you should be able to access

http://ServerName:8180/cas/login and login with your regular account & password

You should also be able to use the ssl connection:

https://ServerName:8443/cas/login

To reset the test, you have to remove the cookies from your browser (Options/clear my traces, uncheck everything but the stuff related to cookies).

Well, if you successfully accessed the pages above, you can now configure your various applications to use the CAS server.

## CASification of the Applications

Now you have a working CAS server. You'll have to integrate CAS to your various applications. It can be done more or less simply, depending on the application and the programming language and/or environment used for these applications.

However, there are some common things to take care of, especially about your cas server's certificate. Most applications require having this certificate trusted/approved/recognized. That means you must declare this certificate in some way.
Tomcat webapps & certificate.

For webapps applications such as jira or confluence, you need to add your cas server's certificate to the certificate repository of the jvm used by your tomcat server.

To do so, first find the JVM used by tomcat. For debian based systems it may be the global jvm in `/usr/lib/jvm`. For jira & confluence, the jvm is embedded with the packages. If your jira is in `/home/jira`, the jvm used might also be under this directory.

```bash
keytool -import -alias cas -file YourCASServerCertificate.pem -keystore $JAVA_HOME/jre/lib/security/cacerts
```

The default password for the jre certificate store is "changeit" ...

### Confluence & Jira

These 2 softwares are provided by the same company and use the same technologies. So their configuration to use CAS is almost identical. You could use 2 different configurations:
http://www.soulwing.org/
http://www.ja-sig.org/wiki/display/CASC/Configuring+Confluence+with+JASIG+CAS+Client+for+Java+3.1

The configuration used in my company is based on the last one. Actually, some modifications were needed for the Single Sign On AND the Single Sign Out to work correctly. Here are the configurations...

#### Confluence

First, add the CAS server's certificate in the jvm keystore as described above.

Go to the directory $CONFLUENCE_HOME/confluence/WEB-INF

Edit the file classes/seraph-config.xml and modify the value of "parameter" for login.url, link.login.url, logout.url to use the cas server. Replace the authenticator class with the jasig cas. You should finally have something like this:

```xml
<security-config>
 <parameters>
  <init-param>
   <param-name>login.url</param-name>
   <param-value>https://deb-cas:8443/cas/login?service=${originalurl}</param-value>
  </init-param>
  <init-param>
   <param-name>link.login.url</param-name>
   <param-value>https://deb-cas:8443/cas/login?service=http://192.168.0.234:18080/</param-value>
  </init-param>
  <init-param>
   <param-name>logout.url</param-name>
   <param-value>https://deb-cas:8443/cas/logout</param-value>
  </init-param>
  <init-param>
   <param-name>cookie.encoding</param-name>
   <param-value>cNf</param-value>
  </init-param>
  <init-param>
   <param-name>login.cookie.key</param-name>
   <param-value>seraph.confluence</param-value>
  </init-param>

  <\!--only basic authentication available-->
  <init-param>
   <param-name>authentication.type</param-name>
   <param-value>os_authType</param-value>
  </init-param>
 </parameters>

 <rolemapper class="com.atlassian.confluence.security.ConfluenceRoleMapper"/>
 <controller class="com.atlassian.confluence.setup.seraph.ConfluenceSecurityController"/>
 <authenticator class="org.jasig.cas.client.integration.atlassian.ConfluenceCasAuthenticator"/>

 <services>
  <service class="com.atlassian.seraph.service.PathService">
   <init-param>
    <param-name>config.file</param-name>
    <param-value>seraph-paths.xml</param-value>
   </init-param>
  </service>
 </services>

 <interceptors>
  <interceptor name="login-logger" class="com.atlassian.confluence.user.ConfluenceLoginInterceptor"/>
 </interceptors>
</security-config>
```

Next, edit the file web.xml and add the following lines after the 2 context-param tokens:

```xml
<!-- CAS Java Client -->
    <filter>
       <filter-name>CAS Single Sign Out Filter</filter-name>
       <filter-class>org.jasig.cas.client.session.SingleSignOutFilter</filter-class>
    </filter>

    <filter>
        <filter-name>CAS Authentication Filter</filter-name>
        <filter-class>org.jasig.cas.client.authentication.AuthenticationFilter</filter-class>
        <init-param>
                <param-name>casServerLoginUrl</param-name>
                <param-value>https://deb-cas:8443/cas/login</param-value>
        </init-param>
        <init-param>
                <param-name>validateUrl</param-name>
                <param-value>https://deb-cas:8443/cas/serviceValidate</param-value>
        </init-param>
        <init-param>
                <param-name>serverName</param-name>
                <param-value>http://192.168.0.234:18080</param-value>
        </init-param>
    </filter>

    <filter>
        <filter-name>CasValidationFilter</filter-name>
        <filter-class>org.jasig.cas.client.validation.Cas20ProxyReceivingTicketValidationFilter</filter-class>
        <init-param>
            <param-name>casServerUrlPrefix</param-name>
            <param-value>https://deb-cas:8443/cas</param-value>
        </init-param>
        <init-param>
            <param-name>serverName</param-name>
            <param-value>http://192.168.0.234:18080</param-value>
        </init-param>
        <init-param>
            <param-name>redirectAfterValidation</param-name>
            <param-value>true</param-value>
        </init-param>
    </filter>

    <filter>
        <filter-name>CAS HttpServletRequest Wrapper Filter</filter-name>
        <filter-class>org.jasig.cas.client.util.HttpServletRequestWrapperFilter</filter-class>
    </filter>

    <filter-mapping>
        <filter-name>CAS Single Sign Out Filter</filter-name>
        <url-pattern>/*</url-pattern>
    </filter-mapping>

    <filter-mapping>
      <filter-name>CAS Authentication Filter</filter-name>
      <url-pattern>/*</url-pattern>
    </filter-mapping>

    <filter-mapping>
        <filter-name>CasValidationFilter</filter-name>
        <url-pattern>/*</url-pattern>
    </filter-mapping>

    <filter-mapping>
        <filter-name>CAS HttpServletRequest Wrapper Filter</filter-name>
        <url-pattern>/*</url-pattern>
    </filter-mapping>

    <listener>
       <listener-class>org.jasig.cas.client.session.SingleSignOutHttpSessionListener</listener-class>
    </listener>

    <!-- End Of CAS Configuration -->
```

Finally, locate and edit the file classes/xwork.xml and modify the redirect parameter in the action.logout item. If you can't locate the file xwork.xml, just unpack it from lib/confluence-\*.jar. You can unpack it with the unzip utility...

You should have something like this:

```xml
         <action name="logout" class="com.atlassian.confluence.user.actions.LogoutAction">
            <interceptor-ref name="defaultStack"/>
            <result name="error" type="velocity">/logout.vm</result>
            <result name="success" type="redirect">https://deb-cas:8443/cas/logout</result>
        </action>
```

Now, you can restart your confluence server. You should be able to login through CAS SSO.

#### Jira

Just proceed as for confluence. In classes/seraph-config.xml, replace Confluence by Jira in the class name org.jasig.cas.client.integration.atlassian.ConfluenceCasAuthenticator
