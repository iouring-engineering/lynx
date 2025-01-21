# Table of Contents

- [Middleware](#id-1)
  - [Configuration](#id-1-1)
  - [Services](#id-1-2)
  - [Database](#id-1-3)
  - [Dynamic Link](#id-1-4)
- [Android](#id-2)
- [IOS](#id-3)
- [Web](#id-4)
- [How it works?](#id-5)

## <center id="id-1" style="font-family:Helvetica ;color:red">Middleware</center>

Here Middleware refers to REST API services which are all used for dynamic linking feature, we used Golang to build API's.

## <span id="id-1-1" style="font-family:Helvetica ;color:lightblue">Configuration</span>

- **Config.yaml**: this file has all the required configuration for the services as well as mobile devices, including database configurations, log file configurations default configuration Open graph(OG meta tags) also android and ios package details.

## <span style="font-family:Helvetica ;color:lightblue">How to run?</span>

- `go run src/*.go -config <config-path>` is the command to run your go services.
- `go build src/*.go -o <output-path>` to generate executable file, `./<executable-file> -config <config-path>` to run the executable file with config.

## <span id="id-1-2" style="font-family:Helvetica ;color:lightblue">Services</span>

- **/.well-known/apple-app-site-association**: this API will be using by APN's server to validate the domain and the package name, as apple server validates randomly we have to make sure that this service should be always up, below is the sample response format which we have to return in response as json.
Refer more from [here][ios-app-verify]

```json
{
  "applinks": {
    "apps": [],
    "details": [
      {
        "appID": "5NG4LBGVER.com.nxtrad.cubeplus",
        "paths": [
          "NOT /_/*",
          "/*"
        ]
      }
    ]
  }
}
```

- **/.well-known/assetlinks.json**: As we do for IOS mobile, same setup is there for android as well, but URI path is different and response format also different,
android mostly validates when we are installing the app itself.
Refer more from [here][android-app-verify].

```json
[
  {
    "relation": [
      "delegate_permission/common.handle_all_urls"
    ],
    "target": {
      "namespace": "android_app",
      "package_name": "com.nxtrad.cubeplus",
      "sha256_cert_fingerprints": [
        "06:67:30:70:4A:40:50:00:CE:43:50:07:55:FB:D1:A1:2A:EC:C5:E4:14:CC:C1:4E:E5:19:F3:8B:58:D5:F0:A7"
      ]
    }
  }
]
```

- **/create**: this API is used to create short link, below json is the expected format of request data will be considered in this API.Its upto developers to handle authentication or authorization for this service, it is recommended to use some auth for this API to avoid that anyone can create the short link.

```json
{
  "android": {
    "fbl": "string" // custom fallback link if the short link opened from the android devices, leave this as empty so that configured behaviour will be handled in android
  },
  "data": any, //any type of data accepted in this field, the same we can get from resolver API's
  "expiry": "string", //2024-01-22 13:08:11 is the expected format
  "ios": {
    "fbl": "string" // custom fallback link if the short link opened from the IOS devices,leave this as empty so that configured behaviour will be handled in ios
  },
  "social": {
    "description": "string",
    "imgUrl": "string",
    "shortIcon": "string",
    "title": "string"
  }, // Open graph fields to show custom fields when we are sharing in social media platforms
  "webUrl": "string" // web url which has to resolved when users clicking on the short link
}
```

- **/{shortcode}**: this API handles navigation based on the different devices, for example if we want bring user to playstore if user doesn't installed the app on their mobile, we just have to configure proper playstore url and update `app-config->android->behavior` as `appsearch` if you want to navigate to any custom webpage update `app-config->android->behavior` as `custom` and change url in `android->android-default-web-url`, similar configuration applies for IOS also, update configurations in
`app-config->ios->behavior` and `app-config->ios-default-web-url`.
- **/data/{shortcode}**: If we want to get the actual data while we used in create API this API can be used, below is the sample data for this API, input field has the actual input for **/create** API.

```json
{
  "s": "ok",
  "d": {
    "input": {
      "type": "ekyc-referral"
    },
    "shortcode": "ebsWWY"
  }
}
```

## <span id="id-1-3" style="font-family:Helvetica ;color:lightblue">Database</span>

We use MySQL as the database to store and retrieve data for the short links.

## <center id="id-2" style="font-family:Helvetica ;color:red">Android</center>

1) In the `AndroidManifest.xml` file located in `android/app/src/main` directory, add an intent filter inside the `<activity>` tag:

![image](docs/android-manifest.png)

**Note** : Ensure to include the `<intent-filter>` configuration within the `AndroidManifest.xml` for all flavors and whitelabel configurations if applicable.

2) Update the `android:launchMode` attribute in the `AndroidManifest.xml` file to `singleTask` to ensure that when a deep link is triggered, the existing app is brought to the foreground instead of creating a new instance. you can learn more about `android:launchMode` [here](<https://medium.com/android-news/android-activity-launch-mode-e0df1aa72242>)

3) Hosting `assetlinks.json` file: Host an `assetlinks.json` file in using a web server with a domain that you own. This file tells the mobile browser which Android application to open instead of the browser. To create the file, get the package name of your app and the sha256 fingerprint of the signing key you will be using to build the APK

  - Locate the package name in `AndroidManifest.xml`, the package property under `<manifest> `tag. Package names are usually in the format of com.example.*.
  - Retrieve the SHA256 fingerprint of the signing key used to build the APK. by using the following command in terminal / command line 

    `keytool -list -v -keystore <path-to-keystore>`

## <center id="id-3" style="font-family:Helvetica ;color:red">IOS</center>
1) Updaing the link in the associated domain capability of iOS application
  - Open the iOS `Runner` app in xcode
  - Click the top-level `Runner`.
  - Click `Signing & Capabilities`.
  - Click `+` Capability to add a new domain (if there isn't one already created. otherwise skip this and the next step).
  - Click `Associated Domains`.
  - In the `Associated Domains` section, click `+`.
  - Enter `applinks:dev-lynx.iourin.in`.

  **Note**: Ensure that the Associated Domain configuration is included within the Runner for all flavors and whitelabel configurations, if applicable.

  2) Hosting `apple-app-site-association` file : You need to host an `apple-app-site-association` file in the web domain. This file tells the mobile browser which iOS application to open instead of the browser. To create the file, get the app ID of the Flutter app.
  - App id: Locate the `bundle ID` in the Xcode project and locate the `team ID` in the [Apple developer account](https://idmsa.apple.com/IDMSWebAuth/signin?appIdKey=891bd3417a7776362562d2197f89480a8547b108fd934911bcbea0110d07f757&path=%2Faccount%2F&rv=1).
  - for example it will be like "S8QB4VV633.com.iouring.nxtrad" where S8QB4VV633 is team ID and com.iouring.nxtrad is bundle id.

**Note**: It might take up to 24 hours before Apple’s Content Delivery Network (CDN) requests the apple-app-site-association (AASA) file from your web domain. The universal link won’t work until the CDN requests the file. 


## <center id="id-4" style="font-family:Helvetica ;color:red">Web</center>

If we are opening the short link in a web browser from any desktop or laptop link will be redirected to the location url which we gave in `webUrl` field when we created the short link, along with the `shortcode` as query param additionally if we are passing json as the input in `data` field when we create short link those json params also will be appended as query params to the web url location when redirecting.

## <center id="id-5" style="font-family:Helvetica ;color:red">How it works?</center>

First and mandatory step to implement dynamic link feature is that our deep link scheme should starts with `http` or `https`, reason is that if someone who doesn't installed the target application android or ios will launch that short link into a default browser only from there our service can take over the control.

Example,

Short link: <https://dev-lynx.iouring.com/PnRj8>

1) If we click above link on android when we have installed the app already with proper configuration for the same url, it will launch the app directly without making any http request to our service.
2) If we don't installed that app already, android will launch default browser with the same short link.
3) In that case our service will be called and with `user-agent` header from the request we will identify which device requested for the link, based on that and configuration for the devices API will send the response.
4) Same implementation for IOS as well, but this feature which is called as universal links supports only from IOS version 9 or greater, for android greater than or equal to android 6.

```mermaid
---
title: Flow diagram
---
flowchart TD
    A[Dynamic Link] --> Android
    A[Dynamic Link] --> IOS
    A[Dynamic Link] --> Desktop
    
    
    Android --> |If| AndroidInstalled[Installed]
    Android --> |If| AndroidNotInstalled[Not Installed]
    AndroidInstalled --> AndroidApp[Opens android app]
    AndroidNotInstalled --> |If Configured </br>for install| Playstore[Opens Play store</br> for installation]
    AndroidNotInstalled --> |If Configured </br> for webpage| AndroidBrowser[Opens a webpage </br>in browser]
    AndroidApp --> AndroidProcess[Reads data from link</br> process it accordingly]
    Playstore --> |After successful</br> installation |AndroidProcess[Reads data from link</br> process it accordingly]

    IOS --> |If| IosInstalled[Installed]
    IOS --> |If| IosNotInstalled[Not Installed]
    IosInstalled --> IosApp[Opens IOS app]
    IosNotInstalled --> |If Configured </br> for install| AppStore[Opens Appstore</br> for installation]
    IosNotInstalled --> |If Configured </br> for webpage| IosBrowser[Opens Appstore</br> for installation]
    IosApp --> IosProcess[Reads data from link</br> process it accordingly]
    AppStore --> |After successful</br> installation |IosProcess[Reads data from link</br> process it accordingly]

    Desktop --> |On Clicking| Browser[Opens a webpage </br>in browser]
```

[android-app-verify]:https://developer.android.com/training/app-links/verify-android-applinks
[ios-app-verify]:https://developer.apple.com/documentation/xcode/supporting-associated-domains

## <center id="id-4" style="font-family:Helvetica ;color:red">Flutter implementaion</center>

- To retrieve the link that opened the app, we utilize a Flutter package named [Unilinks](https://pub.dev/packages/uni_links)

- There are three scenarios that need to be handled:

  1) **The application is in a terminated state**
    - In this case, we will utilize the final `uri = await getInitialUri`() method from the `uni_links` package
  2) **The app is already open and in the background state**

      ![image](docs/handle-incoming-links.png)
  3) The app is not installed, and the link will direct the user to the respective app stores for installation
  - iOS: For iOS applications, before navigating to the `Apple App Store`, we will redirect the user to a webpage. On this webpage, we will copy the referral/deeplink information from the link to the device's clipboard. Subsequently, after the user installs the app and opens it, we will access the clipboard information using the following method.

    `final data = await Clipboard.getData('text/plain');`
    
    **Note** : Ensure that this method is called only once in the entire lifecycle of the app.
  
  - Android: For Android, we utilize the Android PlayInstallReferrer API to retrieve the referral/deeplink information.

    Use the provided link to set up the [PlayInstallReferrer API](https://developer.android.com/google/play/installreferrer/library)

    Utilize Flutter's `method channel` to establish communication with the PlayInstallReferrer API, enabling it to return the referral/deeplink information.

    To learn more about Flutter method channels and their implementation, refer to [this article](https://medium.com/@iiharish97ii/method-channel-in-flutter-bridging-native-code-and-flutter-with-two-way-communication-788d1e91c8c1)
