using UnityEngine;
using UnityEngine.UI;
using System.Collections;

public class ProgressBar : MonoBehaviour {

    const int maxWidth = 512;
    const int height = 20;

    [SerializeField]
    Image foreground;

    [SerializeField]
    public float progress;
	
	// Update is called once per frame
	void Update () {
        progress = Mathf.Min(Mathf.Max(progress, 0), 1);
        foreground.rectTransform.localScale = new Vector3(progress, 1);
	}
}
