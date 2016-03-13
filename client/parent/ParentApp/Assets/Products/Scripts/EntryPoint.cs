using UnityEngine;
using UnityEngine.UI;
using System;
using System.IO;
using System.Collections.Generic;

public class EntryPoint : MonoBehaviour
{

    const string baseURL = "http://localhost:8080/";

    const int softLimit = 100;
    const int hardLimit = 200;

    const string videotag = "20160312";

    [SerializeField]
    Button playBtn;

    [SerializeField]
    RawImage image;

    [SerializeField]
    ProgressBar progressbar;

    Dictionary<int, Texture2D> tex2d = new Dictionary<int, Texture2D>();

    int maxCount;
    int frame;
    bool isPlaying;

    double beforeUpdateTime;

    // Use this for initialization
    void Start()
    {
        isPlaying = false;
        frame = -1;			
        playBtn.onClick.AddListener(() =>
            {
                if (tex2d != null && tex2d.Count > 0)
                {
                    frame = 0;
                    isPlaying = true;
                }
            });
        
        StartCoroutine(GetImageCount(() =>
                {
                    for (var i = 0; i < maxCount; i++)
                    {
                        if (!tex2d.ContainsKey(i))
                        {
                            StartCoroutine(DownloadImage(i));
                        }
                    }
                    GotoFrame(0, true);
                }));
    }
	
    // Update is called once per frame
    void Update()
    {
        if (!isPlaying || frame < 0 || tex2d == null || frame >= maxCount)
        {
            return;
        }
            
        if (frame == 0 || Time.time - beforeUpdateTime >= 0.0333)
        {
            GotoFrame(frame);
            frame++;
            beforeUpdateTime = Time.time;
        }
    }

    IEnumerator<WWW> GetImageCount(Action callback)
    {
        var www = new WWW(baseURL + "/api/imagecount?tag=" + videotag);
        yield return www;
        maxCount = int.Parse(www.text);
        callback();
    }

    IEnumerator<WWW> DownloadImage(int i)
    {               		
        var fname = string.Format("{0}/{1}_{2:D6}.jpg", Application.persistentDataPath, videotag, i + 1);
        var url = string.Format("{0}image/{1}/{2:D6}.jpg", baseURL, videotag, i + 1);
        if (File.Exists(fname))
        {
            var bytes = File.ReadAllBytes(fname);
            var tex = new Texture2D(404, 380);
            tex.LoadImage(bytes);
            AddTex(i, tex);
        }
        else
        {

            Debug.LogFormat("Download {0}...", url);
            var www = new WWW(url);
            yield return www;
            AddTex(i, www.textureNonReadable);

            File.WriteAllBytes(fname, www.bytes);
        }
    }

    void AddTex(int i, Texture2D tex)
    {
        if (tex2d.Count > hardLimit)
        {
            while (tex2d.Count > softLimit)
            {
                tex2d.Remove((new List<int>(tex2d.Keys))[0]);
            }
        }
        if (!tex2d.ContainsKey(i))
        {
            tex2d.Add(i, tex);
        }
    }

    void GotoFrame(int to, bool stop = false)
    {
        if (to < 0 || to >= maxCount)
        {
            return;
        }

        Debug.LogFormat("At {0}", to);

        StartCoroutine(DownloadImage(to));
        if (tex2d.ContainsKey(to))
        {
            image.texture = tex2d[to];
            progressbar.progress = (float) to / maxCount;
        }          

        isPlaying = isPlaying && !stop;
    }
}
