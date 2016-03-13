using UnityEngine;
using UnityEngine.UI;
using System.Collections.Generic;

public class ProgressBar : MonoBehaviour
{

    public enum HighlightType
    {
        
    }

    public class Highlight
    {
        public HighlightType Type { get; private set; }

        public float Location { get; private set; }

        public Highlight(HighlightType type, float location)
        {
            Type = type;
            Location = location;
        }
    }

    [SerializeField]
    Image foreground;

    [SerializeField]
    public float progress;

    private List<Highlight> highlights;

    public ProgressBar()
    {
        highlights = new List<Highlight>();
    }
	
    void Update()
    {
        progress = Mathf.Min(Mathf.Max(progress, 0), 1);
        foreground.rectTransform.localScale = new Vector3(progress, 1);
    }
}
