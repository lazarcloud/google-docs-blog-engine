DESCRIPTION Creating a blog can be a fun and rewarding project, but what if you could simplify the process by using Google Docs as your content management system?

![][image1]

Creating a blog can be a **fun and rewarding** project, but what if you could simplify the process by using **Google Docs** as your **content management system**? In this guide, I'll walk you through how to set up your personal blog using Google Docs, Google Drive, and a bit of web development with Golang and Docker. This method streamlines your content creation process, allowing you to write your posts directly in Google Docs while the setup automates the conversion to blog posts.

\#\#\# Prerequisites:

Before we start, make sure you have the following:

- **A Google account**: You’ll need this for Google Docs and Google Drive.  
- **A hosting platform**: This could be a cloud service like Google Cloud, AWS, or even a personal server that can run Docker containers or virtual machines.  
- **Basic web development skills**: The heavy lifting is already done by using templates and repositories, but some customization may require knowledge of HTML, CSS, and Golang.   
- **Creative content**: Have some ideas you want to share\! This guide is about how to technically set up the blog, but the magic comes from your content.


\#\# 1\. Set up a google service account with json credentials and enable Drive API.

The first step is to allow your application to access your Google Drive, where the blog posts will be stored in Google Docs. This requires creating a Google service account and enabling the Drive API.

\#\#\#\# Steps:  
1\. Go to [Google Cloud Console](https://console.cloud.google.com):  
\- Log in with your Google account.  
\- Create a new project by navigating to **Projects** in the upper left corner and selecting **New Project**.  
\- Give your project a name like "Google Docs Blog" and click **Create**.  
![][image2]

2\. **Create a Service Account**:  
\- Once your project is created, click on the **hamburger menu** in the top left corner.  
\- Navigate to **IAM & Admin** \> **Service Accounts**.  
\- Click **Create Service Account** and fill in the details (e.g., name it "Blog Service Account").  
\- For the role, select **Viewer** to allow the service account to manage your Google Drive.![][image3]![][image4]  
3\. **Add JSON Credentials**:  
\- After creating the service account, click on **Keys** \> **Add Key** \> **Create New Key** and select JSON.  
\- Download the credentials file and keep it safe. You’ll use this to authenticate your application.   
![][image5]

4\. **Enable the Drive API**:  
\- Go to the **API & Services** section in Google Cloud Console.  
\- Click on **Enable APIs and Services** and search for the **Google Drive API**.  
\- Click **Enable** to activate it for your project.![][image6]  
This service account will now act as an intermediary to access your Google Docs for blog content.

\#\# 2\. Set up a new folder in your Google Drive and share it.

To streamline the process, all your blog posts will be stored in a specific folder in your Google Drive. Here's how to set it up.

1\. **Create a Folder**:  
\- Go to [Google Drive](https://drive.google.com). \- Create a new folder, name it something like "Blog Posts", and store all your Google Docs here. 

![][image7]

2\. **Share the Folder with Read Access**:  
\- Right-click the folder and click **Share**.  
\- Set the permission to **Anyone with the link can view** or add the email of your service account. This ensures that your blog app can access the contents without needing manual approvals every time.   
![][image8]  
Only allow read access.  
3\. **Configure Your App with the Folder ID**:  
\- The folder ID is the string of characters after \`folders/\` in the URL when you open the folder.  
\- You'll use this ID to programmatically fetch your blog posts later.

An example of a .env file using all the credentials can be found [here](https://github.com/lazarcloud/google-docs-blog-quickstart/blob/main/example.env).

\#\#\#\# Tip: Make Docs Pages Seamless  
\- To make your blog posts look cleaner, open each Google Doc and go to  **File** \> **Page Setup** and select **Pageless** mode. This will make your blog posts flow more smoothly, without traditional page breaks.  
.![][image9]  
For an example of a blog check this post’s corresponding [docs file](https://docs.google.com/document/d/1XoR0FMwOJjVyYaiQ_Jm7T5nZzCP7d98X1eTTBUivZ3s/edit?usp=sharing).

---

\#\#\# 3\. Set up the Starter Repository  
Now that your Google Drive and service account are ready, it’s time to set up the blog engine. We’ve prepared a template repository that does most of the work for you, allowing you to focus on customization.  
\#\#\#\# Steps:

1\. **Clone the Template Repository**: \- Go to [Google Docs Blog Quickstart](https://github.com/lazarcloud/google-docs-blog-quickstart) on GitHub.  
\- Clone the repository locally.  
\`\`\`bash  
git clone [https://github.com/lazarcloud/google-docs-blog-quickstart.git](https://github.com/lazarcloud/google-docs-blog-quickstart.git)  
\`\`\`  
\- This repository contains everything you need to get started, including a blog engine built with Golang and a front-end powered by Astro.js.

2\. sudo docker system prune: \- Inside the repository, you’ll find two key components: 1\. **app directory**: \- This is where the front-end blog template lives. The template uses [Astro.js](https://astro.build/), which allows you to create fast, static websites.  
\- You can customize the design using basic HTML and CSS.  
\- Keep in mind that any changes in \`./app/src/content/blog\` and \`./app/src/public/images\` folders may be overwritten when new blog posts are generated.  
2\. **main.go**: \- This file runs my [Google Docs Blogs Engine](https://github.com/lazarcloud/google-docs-blog-engine) which connects to your Google Docs folder using the Drive API and converts each Google Doc into a blog post.   
\- You can extend the functionality by using my [Golang Web Server Framework](https://github.com/lazarcloud/google-docs-blog-engine) that is used to serve the static files generated by Astro.js.  
\- For development, uncomment the first lines in \`main.go\` and create a \`.env\` file with your Google variables (credentials, folder ID, etc.). Here’s an [example](https://github.com/lazarcloud/google-docs-blog-quickstart/blob/main/example.env).  
3\. **Setting up the Dockerfile**: \- The repository also contains a pre-configured \`Dockerfile\` to easily deploy your blog as a container. \- If you plan on using a backup function to store your blog posts, make sure to set up appropriate volumes to store data between restarts.

---

\#\#\# Bonus: Customizing the Blog Design If you want to go beyond the basics, here are a few ways to customize your blog:  
\- **Custom CSS**: Modify the \`app\` directory’s stylesheets to give your blog a unique look and feel.  
\- **Change the Layout**: Update the Astro.js templates to reorganize content, create custom sections, or add widgets like social media links, recent posts, etc.  
\- **Additional Features**: Use Golang to integrate new features, such as a commenting system, blog post tagging, or even a subscription feature.

---

\#\#\# Final Thoughts

Setting up a blog using Google Docs can be a powerful way to streamline your content creation. Not only do you get the ease of writing in Google Docs, but with a little setup, your blog posts can be automatically converted into a professional-looking website. If you're a creator who values simplicity and automation, this method can save you time and let you focus on what matters—sharing your ideas with the world. 

[image1]: https://lh7-rt.googleusercontent.com/docsz/AD_4nXf8X26joBWm4mm3rwbw6FE0h1ziAQlvaZWGEdgfRvnicKcvnSP6ddY_6O7gjTTkTJJSjT-k-pdldYwaa5Eyxt87rQSK7vo-CIKxcbLvuwY_N7IDJGwZGcSXz-BtNw_0XJqLZ7rOfZ9vURfjwvKyy6H5z5vEtg?key=KHU3ypSo5_0z_rNIqMmIZg

[image2]: https://lh7-rt.googleusercontent.com/docsz/AD_4nXeJZNHisWWfBllXeAg5qZOlY2hvftRYalnJir5zDGYl0aKPuL5VF2jclljWk2Hn9AQWQ1_vHDB18KyNZc69KAs0yyuv7bopX1sEEBqi3PITPHIj9RswZTEOVwQ3JVPIoxZ1pZVbk6Fn0UisXy5NSSsWcC-quA?key=KHU3ypSo5_0z_rNIqMmIZg

[image3]: https://lh7-rt.googleusercontent.com/docsz/AD_4nXfM-9daSnjHzL8qBxSDcGpZBGqMm0zzHEIplCU6k3W6ouDHzzqkahIxUqGWWAGy41RvUBQO2WzrC6RvbJNvBVm0J4MhbAsTbpGoUlgBnd3DB6tuvYwyxdnZT6hLQ4XTWdmvY_cVRUn5L2TYBlWCvF9te5MjwA?key=KHU3ypSo5_0z_rNIqMmIZg

[image4]: https://lh7-rt.googleusercontent.com/docsz/AD_4nXdefwyI11rCEtIacpd8LzHlbqn-FtVik-v7nFBx_0nJK1zdhWcW_CNknZCBuCQI_03VQJjXer-7YoC6zq_v0YiYolXv9T2r25jb3ksSTS9k9TBbTTDrBXrmNFZBPYujqhS9HCFrPGg-44GZfIVj8sDFRhvUbg?key=KHU3ypSo5_0z_rNIqMmIZg

[image5]: https://lh7-rt.googleusercontent.com/docsz/AD_4nXdArzxLZKteRR_-MnY91eIO_zYUYKMx7S4ibV12GHL37wUhbPeAjV49zBb7khbEHC7T3vxbk06WvH57O0-PmIUXjMoQLgGHRsFEecN7pLDcIj5OWL2b0kIi782aAof4qqGLEr7Pv1voDgtT9NhLmatvVxy8?key=KHU3ypSo5_0z_rNIqMmIZg

[image6]: https://lh7-rt.googleusercontent.com/docsz/AD_4nXcY26XQLqEINd242_3axXiN7JhPbTLZwwU-X2WW3kMXvbH6o53B27h_CRD8JWV2GKv5lpm10qFRHl4JerXesOm4Gixb2DLIIh86EKIR_RVMAUl1bBoF1oT6CJeX0AWL9umSwF-y1Zph_w-dR9oSB2Yu9DUfVA?key=KHU3ypSo5_0z_rNIqMmIZg

[image7]: https://lh7-rt.googleusercontent.com/docsz/AD_4nXc5cceZnjhTpA2igbJ3PMZcuBO-qXq76u8XA5fi6NfLaaVLP7RLw-qKJpm_M5Vhogr8XKp4r3MKvLWkRAMhufJTpitNrU7o2RkUaKyo2ElLZuQTjs1Z8mdz929cg-QR7fZQvFVbbRm8LSihMIZ5pjNWolVQ?key=KHU3ypSo5_0z_rNIqMmIZg

[image8]: https://lh7-rt.googleusercontent.com/docsz/AD_4nXd4FKYRqxwGXoBs4nIadohzCpPB94JKeWVegB7SNwsNlrYaukO-Z77o-T-Kk2E9XsNNmNybezNwZzw3TbILD2cilbI5UV0LnEKSw90hhSmJ2j8BlX9ntB_tg89bT_PjujuC9gbGd-DGUh12CvoltnNl7hZRTA?key=KHU3ypSo5_0z_rNIqMmIZg

[image9]: https://lh7-rt.googleusercontent.com/docsz/AD_4nXdNd7whqLBUeAf1L3sEZX72VJ0HiG0AUfZscMxLz9tK70PeRpNptzzBtj-tzsESmGnzvqGXxVxDilzQD9MQNiypv5E8jyD7XQeOMb9Cic1KoDm7K4TXt-qy8CHdS_V8mX5xwWEG79SqaDqM6xHbayjuqkANXQ?key=KHU3ypSo5_0z_rNIqMmIZg