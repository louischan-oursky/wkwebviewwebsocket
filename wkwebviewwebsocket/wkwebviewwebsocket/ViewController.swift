//
//  ViewController.swift
//  wkwebviewwebsocket
//
//  Created by louischan on 10/5/2023.
//

import UIKit
import WebKit

class ViewController: UIViewController {

    var webview: WKWebView!

    override func viewDidLoad() {
        super.viewDidLoad()
        let config = WKWebViewConfiguration()
        self.webview = WKWebView.init(frame: .zero, configuration: config)
        self.webview.isInspectable = true
        self.webview.translatesAutoresizingMaskIntoConstraints = false
        self.view.addSubview(self.webview)
        self.webview.topAnchor.constraint(equalTo: self.view.safeAreaLayoutGuide.topAnchor).isActive = true
        self.webview.bottomAnchor.constraint(equalTo: self.view.safeAreaLayoutGuide.bottomAnchor).isActive = true
        self.webview.leadingAnchor.constraint(equalTo: self.view.safeAreaLayoutGuide.leadingAnchor).isActive = true
        self.webview.trailingAnchor.constraint(equalTo: self.view.safeAreaLayoutGuide.trailingAnchor).isActive = true
        let url = URL(string: "http://172.20.10.5:4000")!
        let request = URLRequest(url: url)
        self.webview.load(request)
    }
}

