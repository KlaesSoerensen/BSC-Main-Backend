namespace ASP.NETCoreWebAPI.Models
{
    public class Preferences
    {
        public int Id { get; set; } // Primary Key
        public required string Key { get; set; } // Attribute - Required, represents the preference key (e.g., "Language")
        public required string ChosenValue { get; set; } // Attribute - Required, the value chosen by the player (e.g., "DK")
        public string[] AvailableValues { get; set; } // Attribute - Array of available options for the preference (e.g., "DA", "NO", "EN", "SE")
        public required Player Player { get; set; } // Foreign Key (Player) - Required, the player this preference is associated with
    }
}
